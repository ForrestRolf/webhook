import Home from "../screen/Home.vue"
import Hook from "../screen/Hook.vue"
import Logs from "../screen/Logs.vue"
import Template from "../screen/Template.vue";
import Email from "../screen/setting/Email.vue";
import Setting from "../screen/Setting.vue";

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
    },
    {
        name: "setting",
        path: "/setting",
        component: Setting,
        meta: {
            menuKey: "setting"
        },
        redirect: {name: "email-setting"},
        children: [
            {
                name: "email-setting",
                path: "email",
                component: Email,
                meta: {
                    menuKey: "setting",
                    subMenuKey: "email-setting"
                }
            }
        ]
    }
]

export default routes
