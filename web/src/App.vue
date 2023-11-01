<script setup>
import {computed, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";

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
</script>

<template>
    <a-layout>
        <a-layout-header>
            <div class="logo"/>
            <a-menu
                v-model:selectedKeys="selectedKeys"
                theme="dark"
                mode="horizontal"
                :style="{ lineHeight: '64px' }"
            >
                <a-menu-item key="home">Hooks</a-menu-item>
                <a-menu-item key="logs">Logs</a-menu-item>
            </a-menu>
        </a-layout-header>
        <a-layout-content>
            <router-view class="viewport"></router-view>
        </a-layout-content>
    </a-layout>
</template>

<style scoped>

</style>
