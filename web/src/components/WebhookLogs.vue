<script setup>
import {computed, ref, watch} from "vue";
import useAxios from "../support/axios.js";

const props = defineProps({
    webhookId: {
        type: String
    }
})
const visible = ref(false)
const {httpGet} = useAxios()

const id = computed(() => {
    return props.webhookId
})
const tags = {
    "OK": "green",
    "warn": "orange",
    "error": "red",
    "info": "blue",
    "debug": "purple",
}

const fetchLogs = () => {
    httpGet(`/logs`).withQuery({id: id.value}).exec().then(data => {
        logs.value = data.payload
    }).catch(e => {

    }).finally(() => {

    })
}
const show = () => {
    visible.value = true
}
const hide = () => {
    visible.value = false
}
const logs = ref([])
const columns = ref([
    {
        title: 'Level',
        dataIndex: 'level',
        key: 'level',
    },
    {
        title: 'Webhook',
        dataIndex: 'webhookName',
        key: 'webhookName',
    },
    {
        title: 'Message',
        dataIndex: 'message',
        key: 'message',
    },
    {
        title: "Date",
        dataIndex: "created",
        key: "created"
    }
])

watch(visible, (n) => {
    visible.value && fetchLogs()
})

defineExpose({show, hide})
</script>

<template>
    <a-drawer
        v-model:open="visible"
        class="custom-class"
        root-class-name="root-class-name"
        title="Logs"
        placement="right"
        width="40%"
    >
        <a-table :dataSource="logs" :columns="columns">
            <template #bodyCell="{ text, record, index, column }">
                <template v-if="column.key === 'level'">
                    <a-tag :color="tags[text] || 'blue'">{{ text }}</a-tag>
                </template>
            </template>
        </a-table>
    </a-drawer>
</template>

<style scoped>

</style>
