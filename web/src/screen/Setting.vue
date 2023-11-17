<script setup>
import {computed, h, ref} from "vue";
import {MailOutlined, CommentOutlined} from "@ant-design/icons-vue";
import {useRoute, useRouter} from "vue-router";

const route = useRoute()
const router = useRouter()
const selectedKeys = computed({
    set(n) {
        router.push({name: n[0]})
    },
    get() {
        return [route.meta.subMenuKey]
    }
})
const menus = ref([
    {
        key: 'email-setting',
        icon: () => h(MailOutlined),
        label: 'Email profile',
        title: 'Email profile',
    },
    {
        key: 'sms-setting',
        icon: () => h(CommentOutlined),
        label: 'SMS profile',
        title: 'SMS profile',
    },
])
</script>

<template>
    <a-row :gutter="[12, 12]">
        <a-col :span="3">
            <a-menu
                v-model:selectedKeys="selectedKeys"
                theme="light"
                mode="vertical"
                :items="menus"
                :style="{ lineHeight: '64px' }"
            >
            </a-menu>
        </a-col>
        <a-col :span="21">
            <router-view></router-view>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
