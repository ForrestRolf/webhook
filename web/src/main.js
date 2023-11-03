import { createApp } from 'vue'
import App from './App.vue'

import Antd from "ant-design-vue"
import "ant-design-vue/dist/reset.css"
import router from "./router"
import './assets/less/style.less'

const app = createApp(App)
app.use(Antd)
app.use(router)

app.mount('#app')