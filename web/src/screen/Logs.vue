<script setup>
import {onMounted, ref, watch} from "vue";
import useAxios from "../support/axios.js";
import {ClearOutlined} from "@ant-design/icons-vue"
import loglevelColor from "../support/log-level.js"
import useMessage from "../support/message.js";
const {httpGet, httpDelete} = useAxios()
const {successMessage, errorMessage} = useMessage()

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
        title: 'Log ID',
        dataIndex: 'logId',
        key: 'logId',
    },
    {
        title: 'Message',
        dataIndex: 'message',
        key: 'message',
        width: "70%",
    },
    {
        title: 'Action driver',
        dataIndex: 'actionDriver',
        key: 'actionDriver',
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

const clearing = ref(false)
const clearLogs = () => {
    clearing.value = true
    httpDelete(`/logs`).exec().then(data => {
        successMessage("Deleted successfully", `Total number of deleted logs: ${data.payload}`).show()
    }).catch(e => {
        errorMessage("Unable to clear logs", e).show()
    }).finally(() => {
        clearing.value = false
    })
}
onMounted(() => {
    setTimeout(fetchLogs, 200)
})
</script>

<template>
    <a-row :gutter="[0, 12]">
        <a-col :span="24" class="txt-rgt">
            <a-button size="small" danger :loading="clearing" @click="clearLogs">
                <template #icon>
                    <ClearOutlined />
                </template>
                Clear last 4 weeks logs
            </a-button>
        </a-col>
        <a-col :span="24">
            <a-table :dataSource="logs" :columns="columns" :pagination="{defaultPageSize: 100}" bordered size="small" :loading="loading">
                <template #bodyCell="{ text, record, index, column }">
                    <template v-if="column.key === 'level'">
                        <a-tag :color="loglevelColor[text] || 'blue'">{{ text.toUpperCase() }}</a-tag>
                    </template>
                </template>
            </a-table>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>