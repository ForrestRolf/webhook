<script setup>
import {onMounted, ref, watch} from "vue";
import useAxios from "../support/axios.js";
import loglevelColor from "../support/log-level.js"
const {httpGet} = useAxios()
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
        title: 'Webhook ID',
        dataIndex: 'webhookId',
        key: 'webhookId',
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
const loading = ref(true)
const fetchLogs = () => {
    loading.value = true
    httpGet(`/logs`).withQuery({limit: 5000}).exec().then(data => {
        logs.value = data.payload
    }).catch(e => {

    }).finally(() => {
        loading.value = false
    })
}
onMounted(() => {
    setTimeout(fetchLogs, 200)
})
</script>

<template>
    <a-table :dataSource="logs" :columns="columns" :pagination="{defaultPageSize: 100}" bordered size="small" :loading="loading">
        <template #bodyCell="{ text, record, index, column }">
            <template v-if="column.key === 'level'">
                <a-tag :color="loglevelColor[text] || 'blue'">{{ text.toUpperCase() }}</a-tag>
            </template>
        </template>
    </a-table>
</template>

<style scoped>

</style>