# OpenFeed

خواندن کانال‌های عمومی تلگرام از طریق یک PWA سبک، حتی زیر فیلترینگ.

## چطور کار می‌کند؟

OpenFeed یک سرور کوچک نوشته‌شده به Go است که صفحه‌ی پیش‌نمایش عمومی
تلگرام (`t.me/s/<channel>`) را نه به‌صورت مستقیم، بلکه از طریق
**Domain Fronting با گوگل‌ترنسلیت** به‌همراه **جعل اثرانگشت TLS (uTLS)**
دریافت می‌کند. این یعنی حتی در شرایطی که خودِ تلگرام فیلتر باشد،
محتوای کانال‌های عمومی همچنان قابل‌خواندن است.

نتیجه به‌صورت JSON تمیز از طریق یک API محلی در اختیار یک PWA (Progressive
Web App) قرار می‌گیرد که مستقیماً در مرورگر گوشی یا کامپیوتر اجرا می‌شود
— بدون نیاز به نصب اپ جداگانه.

## امکانات

- خواندن پست‌های متنی، عکس و ویدیوی کانال‌های عمومی تلگرام
- دانلود واقعی فایل/ویدیو/صدا در صورتی که پیش‌نمایش تلگرام آن را embed
  کرده باشد (نه فقط تصویر بندانگشتی)
- ذخیره‌ی کانال‌های موردعلاقه و نمایش آفلاین از کش مرورگر
- Pull-to-refresh، حالت PWA قابل‌نصب روی صفحه‌ی اصلی گوشی
- طراحی واکنش‌گرا و راست‌به‌چپ (RTL) با پشتیبانی از حالت روشن/تیره

## نصب و اجرا

نیازمند [Go](https://go.dev) نسخه‌ی ۱.۲۵ به بالا.

### نصب معمولی

```bash
git clone https://github.com/yusi20006-max/Openfeed.git
cd Openfeed
go mod download
go build -o openfeed ./cmd/server
./openfeed
```

سپس مرورگر را باز کنید روی:

```text
http://127.0.0.1:8080
```

### نصب روی Termux

مسیر نصب Termux به‌صورت واقعی روی Termux/Android تست و تأیید شده است.
ابتدا مطمئن شوید Go نسخه‌ی ۱.۲۵ یا بالاتر نصب است.

#### دستور کامل نصب از صفر

```bash
cd ~/projects && rm -rf Openfeed && pkg update && pkg install -y git golang && git clone https://github.com/yusi20006-max/Openfeed.git && cd Openfeed && bash scripts/termux-bootstrap.sh && ./openfeed
```

این دستور کل مسیر را انجام می‌دهد:

```text
update packages
    ↓
install Git + Go
    ↓
clone OpenFeed
    ↓
Termux bootstrap
    ↓
download Go dependencies
    ↓
build OpenFeed
    ↓
run OpenFeed
```

> اگر نمی‌خواهید نسخه‌ی قبلی پروژه حذف شود، بخش `rm -rf Openfeed` را اجرا نکنید و فقط وارد repository موجود شوید.

مسیر مرحله‌ای معادل:

```bash
pkg update
pkg install -y git golang
git clone https://github.com/yusi20006-max/Openfeed.git
cd Openfeed
bash scripts/termux-bootstrap.sh
./openfeed
```

اسکریپت Termux قبل از build، dependencyها را با fallback چندمرحله‌ای دریافت می‌کند:

```text
https://goproxy.cn → https://proxy.golang.org → direct
```

این موضوع برای شبکه‌هایی مهم است که `proxy.golang.org` با خطای `403 Forbidden` قابل دسترسی نیست. در چنین شرایطی نباید `go run` را چند بار تکرار کرد؛ ابتدا bootstrap را اجرا کنید تا dependencyها آماده و binary ساخته شود.

اگر mirror اول در دسترس نبود، fallback بعدی امتحان می‌شود. در صورت نیاز می‌توانید مسیر مستقیم را نیز امتحان کنید:

```bash
GOPROXY=direct GOSUMDB=off go mod download
go build -o openfeed ./cmd/server
```

اسکریپت فقط dependencyهای Go را آماده و binary را build می‌کند و هیچ تنظیم runtime یا فایل خارج از پروژه را تغییر نمی‌دهد.

## ساختار پروژه

```text
cmd/server/        نقطه‌ی ورود سرور
internal/telemirror/  کلاینت اصلی: fronting، uTLS، پارس HTML تلگرام
internal/provider/     لایه‌ی انتخاب روش دریافت (پیش‌فرض: telemirror)
internal/parser/       تبدیل داده‌ی داخلی به مدل عمومی API
internal/model/        ساختار JSON که به فرانت‌اند داده می‌شود
internal/api/          هندلرهای HTTP (/api/channel, /api/status, /api/download)
web/                    PWA (HTML/CSS/JS ساده، بدون فریم‌ورک)
scripts/               ابزارهای bootstrap نصب و build
```

## محدودیت‌ها

- فقط کانال‌های **عمومی** تلگرام قابل‌خواندن هستند (نه گروه‌ها یا چت‌های خصوصی)
- دانلود فایل واقعی فقط برای مدیایی ممکن است که پیش‌نمایش تلگرام خودش
  آن را در HTML قرار داده باشد؛ ویدیوهای بزرگ گاهی فقط تصویر بندانگشتی دارند
- سقف حجم دانلود فعلی ۲۰۰ مگابایت است

## مجوز

این پروژه شخصی است؛ فایل مجوز به‌دلخواه بعداً اضافه می‌شود.



---

# Stable Core Policy

OpenFeed is considered the Stable Core of the Yasin ecosystem.

Its responsibilities are:

- Fetch content from supported sources
- Provide a stable local API
- Download and cache media
- Maintain high reliability

New features such as AI processing, scheduling, publishing,
queue management, dashboards, and plugins must be implemented
in downstream projects such as FeedBridge.

OpenFeed will receive only:

- Bug fixes
- Security updates
- Performance improvements
- Compatibility updates

The Termux work in this release is strictly a compatibility and
installation improvement. It must not change Stable Core runtime behavior,
API contracts, fetching semantics, or PWA behavior.
