function getUrlToken() {
    const s = window.location.search;
    let tk = "";
    if (s.length > 1 && s.indexOf("hotspot_tk") > 0) {
        const f = s
            .substring(1)
            .split("&")
            .find((i) => i.indexOf("hotspot_tk=") === 0);
        tk = f ? f.substring(11) : "";
    }
    return tk;
}

function removeUrlHotspotToken(){
    const s = window.location.search;
    let tk = "";
    if (s.length > 1 && s.indexOf("hotspot_tk") > 0) {
        const f = s.substring(1).split("&").filter((i) => i.indexOf("hotspot_tk=") !== 0);
        window.location.search = '?' + f.join("&")
    }
}

export const AuthToken = {
    get() {
        let tk = getUrlToken();
        if (tk) {
            return tk;
        } else {
            return localStorage.hotspotAuthToken;
        }
    },
    set(value) {
        localStorage.hotspotAuthToken = value;
    },
    clean() {
        AuthToken.set("");
        removeUrlHotspotToken()
    },
    exist() {
        return !!AuthToken.get();
    },
};

let router = null;

export const toLoginView = function () {
    router?.replace('/login')
};

export const setRouter = function (r) {
    router = r;
}