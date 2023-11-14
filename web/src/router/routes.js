import Home from "../screen/Home.vue"
import Hook from "../screen/Hook.vue"
import Logs from "../screen/Logs.vue"
import Template from "../screen/Template.vue";

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
    },
    {
        name: "templates",
        path: "/templates",
        component: Template,
        meta: {
            menuKey: "templates"
        }
    }
]

export default routes
