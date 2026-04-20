import { getCookie, deleteCookie } from "https://cdn.jsdelivr.net/gh/jscroot/lib@0.2.8/cookie.js";
import { redirect } from "https://cdn.jsdelivr.net/gh/jscroot/lib@0.2.8/url.js";
import { setInnerText, onClick } from "https://cdn.jsdelivr.net/gh/jscroot/lib@0.2.8/element.js";

// --- Baca sesi dari cookie ---
const userName    = getCookie("user_name");
const userEmail   = getCookie("user_email");
const userPicture = getCookie("user_picture");
const loginTime   = getCookie("user_login_time");

// Proteksi halaman: redirect ke login jika belum punya sesi
if (userName === "") {
  redirect("index.html");
}

// --- Isi elemen halaman ---
setInnerText("user-name",      userName);
setInnerText("user-email",     userEmail);
setInnerText("detail-email",   userEmail);
setInnerText("detail-name",    userName);
setInnerText("detail-session", loginTime || "Sesi ini");

setAvatar("user-avatar", userPicture, userName);

// --- Logout ---
onClick("logout-btn", function () {
  deleteCookie("user_name");
  deleteCookie("user_email");
  deleteCookie("user_picture");
  deleteCookie("user_login_time");
  redirect("index.html");
});

// --- Helpers ---

// Set src gambar; fallback ke avatar inisial jika URL kosong
function setAvatar(id, src, name) {
  const el = document.getElementById(id);
  if (!el) return;
  el.alt = name;
  el.src = src
    ? src
    : `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=1A1A1A&color=fff&size=64`;
}
