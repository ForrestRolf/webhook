import {createRouter, createWebHistory, createWebHashHistory} from "vue-router"
import routes from "./routes.js"

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes
})

export default router
