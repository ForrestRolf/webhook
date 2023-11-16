<script setup>
import {computed, h, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {NodeExpandOutlined, ExceptionOutlined, CodeOutlined, SettingOutlined} from "@ant-design/icons-vue"

const route = useRoute()
const router = useRouter()

const selectedKeys = computed({
    set(n) {
        router.push({name: n[0]})
    },
    get() {
        return [route.meta.menuKey]
    }
})
const menus = ref([
    {
        key: 'home',
        icon: () => h(NodeExpandOutlined),
        label: 'Webhooks',
        title: 'Webhooks',
    },
    {
        key: 'logs',
        icon: () => h(ExceptionOutlined),
        label: 'Logs',
        title: 'Logs',
    },
    {
        key: 'templates',
        icon: () => h(CodeOutlined),
        label: 'Templates',
        title: 'Templates',
    },
    {
        key: 'setting',
        icon: () => h(SettingOutlined),
        label: 'Setting',
        title: 'Setting',
    }
])
</script>

<template>
    <a-layout>
        <a-layout-header>
            <div class="logo"/>
            <a-menu
                v-model:selectedKeys="selectedKeys"
                theme="light"
                mode="horizontal"
                :items="menus"
                :style="{ lineHeight: '64px' }"
            >
            </a-menu>
        </a-layout-header>
        <a-layout-content>
            <router-view class="viewport"></router-view>
        </a-layout-content>
    </a-layout>
</template>

<style scoped>

</style>
