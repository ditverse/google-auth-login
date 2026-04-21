import { setCookieWithExpireDay, getCookie } from "https://cdn.jsdelivr.net/gh/jscroot/lib@0.2.8/cookie.js";
import { redirect } from "https://cdn.jsdelivr.net/gh/jscroot/lib@0.2.8/url.js";
import { setInnerText, show } from "https://cdn.jsdelivr.net/gh/jscroot/lib@0.2.8/element.js";

// Ganti dengan URL backend; kosongkan ("") untuk mode demo tanpa backend
const BACKEND_URL = "https://github.com/ditverse/google-auth-login.git";

// Redirect ke dashboard jika sudah punya sesi
if (getCookie("user_name") !== "") {
  redirect("dashboard.html");
}

// Callback dari Google Identity Services setelah sign-in berhasil
window.handleCredentialResponse = async function (response) {
  const idToken = response.credential;

  if (BACKEND_URL === "") {
    handleDemoMode(idToken);
    return;
  }

  await handleBackendMode(idToken);
};

// --- Mode backend: kirim token ke server untuk diverifikasi ---
async function handleBackendMode(idToken) {
  try {
    const res = await fetch(`${BACKEND_URL}/auth/google`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token: idToken }),
    });

    if (!res.ok) {
      const { message } = await res.json();
      showError(message || "Login gagal. Coba lagi.");
      return;
    }

    const { name, email, picture } = await res.json();
    saveSessionAndRedirect(name, email, picture);
  } catch {
    showError("Koneksi ke server gagal. Pastikan backend berjalan.");
  }
}

// --- Mode demo: parse JWT secara lokal tanpa verifikasi backend ---
function handleDemoMode(idToken) {
  const payload = parseJWTPayload(idToken);
  if (!payload) {
    showError("Gagal membaca token Google.");
    return;
  }
  saveSessionAndRedirect(payload.name, payload.email, payload.picture);
}

// --- Simpan sesi ke cookie lalu redirect ke dashboard ---
function saveSessionAndRedirect(name, email, picture) {
  const loginTime = new Date().toLocaleString("id-ID", {
    dateStyle: "long",
    timeStyle: "short",
  });
  setCookieWithExpireDay("user_name",       name,      1);
  setCookieWithExpireDay("user_email",      email,     1);
  setCookieWithExpireDay("user_picture",    picture,   1);
  setCookieWithExpireDay("user_login_time", loginTime, 1);
  redirect("dashboard.html");
}

function showError(message) {
  setInnerText("error-msg", message);
  show("error-msg");
}

// Parse payload JWT tanpa verifikasi signature (aman hanya untuk demo)
function parseJWTPayload(token) {
  try {
    const base64 = token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(atob(base64));
  } catch {
    return null;
  }
}
