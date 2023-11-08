import axios from 'axios'
import {Notification} from "./message.js";

const instance = axios.create({
    baseURL: import.meta.env.VITE_API_URL || "/api",
    timeout: 1000 * 60,
    withCredentials: false,
})
const beforeRequest = config => {
    config.headers['Content-Type'] = 'application/json'
    config.headers['Accept'] = 'application/json'

    const token = localStorage.getItem('access_token')
    token && (config.headers.Authorization = token)
    return config;
}
const onSuccess = (response) => {
    let code = response.data?.meta?.code || "Invalid Response"
    let message = response.data?.meta?.message || ""
    if (code === "OK") {
        return Promise.resolve(response.data)
    }
    new Notification("error", code, message).show()
    return Promise.reject(message)
}
const onError = (error) => {
    const {response} = error
    if (error.code === "ERR_NETWORK" || error.code === "ERR_FAILED") {
        new Notification("error", error.code, error.message).show()
        return Promise.reject(error.message)
    }
    if (response) {
        const code = error?.response?.data?.code
        if (response.status === 401) {
            location.reload()
        } else if (response.status === 403) {
            new Notification("error", "Access Denied", "").show()
        } else if (response.status === 422) {
            //pass
        } else {
            new Notification("error", response.statusText, error.message).show()
        }
        return Promise.reject(error?.response?.data)
    }
    return Promise.reject(error.message)
}

instance.interceptors.request.use(beforeRequest)
instance.interceptors.response.use(onSuccess, onError)

function Http(instance, url, method) {
    this.instance = instance
    this.url = url
    this.method = method
    this.headers = null
    this.body = null
    this.query = null
    this.pathVariables = null
    this.auth = null
    this.proxy = null
}

Http.prototype.withUrl = function (url) {
    this.url = url
    return this
}
Http.prototype.withHeaders = function (headers) {
    this.headers = headers
    return this
}
Http.prototype.withBody = function (body) {
    this.body = body
    return this
}
Http.prototype.withData = function (data) {
    this.body = data
    return this
}
Http.prototype.withQuery = function (query) {
    this.query = query
    return this
}
Http.prototype.withPathVariables = function (variables) {
    this.pathVariables = variables
    return this
}
Http.prototype.withAuth = function (auth) {
    this.auth = auth
    return this
}
Http.prototype.withProxy = function (proxy) {
    this.proxy = proxy
    return this
}
Http.prototype.exec = function () {
    return new Promise((resolve, reject) => {
        const config = {
            method: this.method,
        }
        let url = this.url
        if (this.pathVariables) {
            for (let k in this.pathVariables) {
                url = url.replace(`:${k}`, this.pathVariables[k])
            }
        }
        config['url'] = url
        if (this.headers) {
            config['headers'] = this.headers
        }
        if (this.query) {
            config['params'] = this.query
        }
        if (this.body) {
            config['data'] = this.body
        }
        if (this.auth) {
            config['auth'] = this.auth
        }
        if (this.proxy) {
            config['proxy'] = this.proxy
        }
        this.instance.request(config).then(resolve).catch(reject)
    })
}

export default function useAxios() {
    const httpGet = (url) => {
        return new Http(instance, url, 'GET')
    }
    const httpPost = (url) => {
        return new Http(instance, url, 'POST')
    }
    const httpPut = (url) => {
        return new Http(instance, url, "PUT")
    }
    const httpDelete = (url) => {
        return new Http(instance, url, "DELETE")
    }

    return {
        httpGet, httpPost, httpPut, httpDelete
    }
}
