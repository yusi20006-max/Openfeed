(function setupMobileMenu(){
  const header = document.querySelector("header");
  const savedWrap = document.getElementById("savedWrap");
  const savedList = document.getElementById("savedList");
  if(!header || !savedWrap || !savedList) return;

  const button = document.createElement("button");
  button.type = "button";
  button.className = "mobile-menu-btn";
  button.setAttribute("aria-label", "باز کردن منو");
  button.setAttribute("aria-expanded", "false");
  button.textContent = "☰";

  const overlay = document.createElement("div");
  overlay.className = "mobile-menu-overlay";
  overlay.setAttribute("aria-hidden", "true");

  const drawer = document.createElement("aside");
  drawer.className = "mobile-menu-drawer";
  drawer.setAttribute("aria-label", "منوی کانال‌ها");
  drawer.innerHTML = `
    <div class="mobile-menu-head">
      <strong>کانال‌ها</strong>
      <button type="button" class="mobile-menu-close" aria-label="بستن منو">×</button>
    </div>
    <div class="mobile-menu-content">
      <div class="mobile-menu-label">کانال‌های ذخیره‌شده</div>
      <div class="mobile-saved-list"></div>
    </div>
  `;

  document.body.appendChild(overlay);
  document.body.appendChild(drawer);
  header.insertBefore(button, header.firstChild);

  const mirror = drawer.querySelector(".mobile-saved-list");
  const closeButton = drawer.querySelector(".mobile-menu-close");

  function syncSaved(){
    mirror.innerHTML = savedList.innerHTML;
    mirror.querySelectorAll(".chip").forEach((chip, index) => {
      chip.onclick = (event) => {
        const original = savedList.querySelectorAll(".chip")[index];
        if(original) original.click();
        event.stopPropagation();
        close();
      };
      const remove = chip.querySelector(".x");
      if(remove) {
        remove.onclick = (event) => {
          const original = savedList.querySelectorAll(".chip")[index];
          const originalRemove = original && original.querySelector(".x");
          if(originalRemove) originalRemove.click();
          event.stopPropagation();
          syncSaved();
        };
      }
    });
  }

  function open(){
    syncSaved();
    drawer.classList.add("open");
    overlay.classList.add("open");
    button.setAttribute("aria-expanded", "true");
    overlay.setAttribute("aria-hidden", "false");
    document.body.classList.add("mobile-menu-open");
  }

  function close(){
    drawer.classList.remove("open");
    overlay.classList.remove("open");
    button.setAttribute("aria-expanded", "false");
    overlay.setAttribute("aria-hidden", "true");
    document.body.classList.remove("mobile-menu-open");
  }

  button.addEventListener("click", open);
  closeButton.addEventListener("click", close);
  overlay.addEventListener("click", close);
  document.addEventListener("keydown", event => {
    if(event.key === "Escape") close();
  });

  const observer = new MutationObserver(syncSaved);
  observer.observe(savedList, { childList: true, subtree: true });
  syncSaved();
})();
