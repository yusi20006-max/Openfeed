(function(){
  const nativeFetch = window.fetch.bind(window);
  const state = Object.create(null);
  const prefix = "/api/channel/";

  function channelFrom(url){
    try{
      const u = new URL(url, location.href);
      return u.pathname.startsWith(prefix) ? decodeURIComponent(u.pathname.slice(prefix.length)) : "";
    }catch(e){ return ""; }
  }

  // The existing app fetches /api/channel/<name>. Upgrade those requests
  // to the 50-post API without changing the existing load flow.
  window.fetch = async function(input, init){
    let url;
    try{ url = new URL(typeof input === "string" ? input : input.url, location.href); }
    catch(e){ return nativeFetch(input, init); }
    const name = channelFrom(url.toString());
    if(!name) return nativeFetch(input, init);
    if(!url.searchParams.has("limit")) url.searchParams.set("limit", "50");
    const response = await nativeFetch(url.toString(), init);
    try{
      state[name] = await response.clone().json();
      window.__openfeedState = state;
    }catch(e){}
    return response;
  };

  function loadMore(){
    const feed = document.getElementById("feed");
    const input = document.getElementById("channel");
    if(!feed || !input) return;
    const name = input.value.trim();
    const data = state[name];
    const old = feed.querySelector(".openfeed-load-more");
    if(old) old.remove();
    if(!data || !data.has_more || !data.next_before) return;

    const wrap = document.createElement("div");
    wrap.className = "openfeed-load-more";
    wrap.style.cssText = "display:flex;justify-content:center;padding:18px 12px 28px";
    const btn = document.createElement("button");
    btn.textContent = "پست‌های بیشتر";
    btn.style.cssText = "border:0;border-radius:12px;padding:12px 24px;font-size:16px;cursor:pointer;min-width:180px;background:#229ed9;color:#fff";
    btn.onclick = async () => {
      btn.disabled = true;
      btn.textContent = "در حال دریافت…";
      try{
        const url = prefix + encodeURIComponent(name) + "?before=" + encodeURIComponent(data.next_before) + "&limit=50";
        const res = await nativeFetch(url);
        if(!res.ok) throw new Error("HTTP " + res.status);
        const older = await res.json();
        const seen = new Set((data.posts || []).map(p => p.id));
        data.posts = (data.posts || []).concat((older.posts || []).filter(p => !seen.has(p.id)));
        data.posts.sort((a,b) => new Date(b.date) - new Date(a.date));
        data.next_before = older.next_before || "";
        data.has_more = !!older.has_more;
        state[name] = data;
        localStorage.setItem("openfeed_cache_" + name, JSON.stringify({channel:data,savedAt:Date.now()}));
        render(data);
      }catch(e){
        btn.disabled = false;
        btn.textContent = "پست‌های بیشتر";
        if(typeof showToast === "function") showToast("دریافت پست‌های بیشتر ناموفق بود");
      }
    };
    wrap.appendChild(btn);
    feed.appendChild(wrap);
  }

  async function renderLatest(){
    const feed = document.getElementById("feed");
    const saved = typeof loadSaved === "function" ? loadSaved() : [];
    if(!feed) return;
    if(!saved.length){
      feed.innerHTML = '<div class="empty">هنوز کانالی ذخیره نشده است</div>';
      return;
    }
    feed.innerHTML = '<div class="loading">در حال دریافت تازه‌ها…</div>';
    const all = [];
    const seen = new Set();
    for(const channel of saved){
      try{
        const res = await nativeFetch(prefix + encodeURIComponent(channel.username) + "?limit=50");
        if(!res.ok) continue;
        const data = await res.json();
        state[channel.username] = data;
        for(const post of (data.posts || [])){
          const key = channel.username + ":" + post.id;
          if(seen.has(key)) continue;
          seen.add(key);
          all.push({...post, source: channel.title || channel.username, sourceUsername: channel.username});
        }
      }catch(e){}
    }
    all.sort((a,b) => new Date(b.date) - new Date(a.date));
    feed.innerHTML = "";
    const head = document.createElement("div");
    head.className = "channel-header";
    const meta = document.createElement("div");
    meta.className = "channel-meta";
    meta.innerHTML = "<h3>تازه‌ها</h3><div class=\"sub\">۵۰ پست آخر کانال‌های ذخیره‌شده</div>";
    head.appendChild(meta);
    feed.appendChild(head);

    all.slice(0,50).forEach(post => {
      const card = document.createElement("div");
      card.className = "post";
      const source = document.createElement("div");
      source.style.cssText = "font-weight:700;margin-bottom:8px;color:#229ed9";
      source.textContent = "@" + post.sourceUsername + " · " + post.source;
      card.appendChild(source);
      const grid = document.createElement("div");
      grid.className = "media-grid";
      (post.media || []).forEach(m => { const el = renderMedia(m, post.id); if(el) grid.appendChild(el); });
      if(grid.childElementCount) card.appendChild(grid);
      if(post.text){ const text = document.createElement("p"); text.innerHTML = post.text; card.appendChild(text); }
      const row = document.createElement("div");
      row.className = "meta";
      row.innerHTML = `<span>${timeAgo(post.date)}</span><span>${post.views || ""}</span>`;
      card.appendChild(row);
      feed.appendChild(card);
    });
    if(!all.length) feed.insertAdjacentHTML("beforeend", '<div class="empty">پست تازه‌ای برای نمایش پیدا نشد</div>');
  }

  function addLatestMenu(){
    const drawer = document.querySelector(".mobile-menu-drawer");
    const content = drawer && drawer.querySelector(".mobile-menu-content");
    if(!content || content.querySelector(".mobile-latest-btn")) return;
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "mobile-latest-btn";
    btn.textContent = "تازه‌ها";
    btn.style.cssText = "width:100%;border:0;border-radius:10px;padding:11px 14px;margin:0 0 12px;text-align:right;font:inherit;cursor:pointer";
    btn.onclick = () => {
      renderLatest();
      const close = drawer.querySelector(".mobile-menu-close");
      if(close) close.click();
    };
    content.insertBefore(btn, content.firstChild);
  }

  const observer = new MutationObserver(() => { loadMore(); addLatestMenu(); });
  observer.observe(document.body, {childList:true, subtree:true});
  addLatestMenu();
})();
