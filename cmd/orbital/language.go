package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type languageKey struct{}

func normalizeLang(lang string) string {
	if lang == "en" {
		return "en"
	}
	return "ar"
}

func requestLang(r *http.Request) string {
	lang := r.Header.Get("X-Orbital-Lang")
	if lang == "" && r.URL.Path == "/api/colony/stream" {
		lang = r.URL.Query().Get("lang")
	}
	return normalizeLang(lang)
}

func withLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, languageKey{}, normalizeLang(lang))
}

func contextLang(ctx context.Context) string {
	lang, _ := ctx.Value(languageKey{}).(string)
	return normalizeLang(lang)
}

func tr(lang, ar, en string) string {
	if lang == "en" {
		return en
	}
	return ar
}

var errorArabic = map[string]string{
	"Stop the colony first":            "وقّف المستعمرة الأول",
	"bad request":                      "طلب غير صالح",
	"unknown level":                    "مهمة غير معروفة",
	"unknown chapter":                  "فصل غير معروف",
	"not in workspace":                 "المهمة مش موجودة في مساحة العمل",
	"forbidden host":                   "عنوان مضيف غير مسموح",
	"missing token":                    "رمز الدخول مفقود أو غير صالح",
	"streaming unavailable":            "البث غير متاح",
	"invalid profile id":               "معرّف اللاعب غير صالح",
	"name and valid language required": "اكتب اسم واختار لغة صالحة",
	"name required":                    "اكتب الاسم",
	"invalid language":                 "لغة غير صالحة",
	"cannot delete the last profile":   "مينفعش تحذف آخر لاعب",
	"invalid or oversized progress":    "بيانات التقدم غير صالحة أو أكبر من الحد المسموح",
	"invalid JSON":                     "بيانات JSON غير صالحة",
}

func errorText(lang, message string) string {
	if arabic, ok := errorArabic[message]; ok {
		return tr(lang, arabic, message)
	}
	return message
}

func requestError(w http.ResponseWriter, r *http.Request, message string, status int) {
	// server-side failures land in the log (make logs) so they can be diagnosed
	if status >= http.StatusInternalServerError {
		log.Printf("%s %s: %d %s", r.Method, r.URL.Path, status, message)
	}
	if strings.HasPrefix(r.URL.Path, "/api/profiles") {
		profileBadRequest(w, r, message, status)
		return
	}
	http.Error(w, errorText(requestLang(r), message), status)
}
