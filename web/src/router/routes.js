import Home from "../screen/Home.vue"
import Hook from "../screen/Hook.vue"
import Logs from "../screen/Logs.vue"

const routes = [
    {
        name: "home",
        path: "",
        component: Home,
        meta: {
            menuKey: "home"
        }
    },
    {
        name: "hooks",
        path: "/hooks",
        component: Hook,
        meta: {
            menuKey: "home"
        }
    },
    {
        name: "logs",
        path: "/logs",
        component: Logs,
        meta: {
            menuKey: "logs"
        }
    }
]

export default routes
