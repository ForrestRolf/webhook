<script setup>
import {computed, ref, watch} from "vue";
import useAxios from "../support/axios.js";
import loglevelColor from "../support/log-level.js"
import LogId from "./LogId.vue";

const props = defineProps({
    webhookId: {
        type: String
    }
})
const visible = ref(false)
const {httpGet} = useAxios()
const loading = ref(false)
const colors = [
    "#4b38b3",
    "#299cdb",
    "#45CB85",
    "#f06548",
    "#ffbe0b",
    "#02a8b5",
    "#00e272",
    "#d568fb",
    "#feb56a",
    "#544fc5",
    "#fe6a35"
]
const pickedColor = ref({})

const id = computed(() => {
    return props.webhookId
})

const fetchLogs = () => {
    loading.value = true
    httpGet(`/logs`).withQuery({id: id.value, limit: 1000}).exec().then(data => {
        logs.value = data.payload
    }).catch(e => {

    }).finally(() => {
        loading.value = false
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
        width: "70%"
    },
    {
        title: "Date",
        dataIndex: "created",
        key: "created"
    }
])

const pickLogIdColor = (id) => {
    if(pickedColor.value[id]) {
        return pickedColor.value[id]
    }
    let c = colors.pop()
    pickedColor.value[id] = c || "#000"
    return c || "#000"
}

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
        width="65%"
    >
        <a-table :dataSource="logs" :columns="columns" :pagination="{defaultPageSize: 20}" size="small" :loading="loading">
            <template #bodyCell="{ text, record, index, column }">
                <template v-if="column.key === 'level'">
                    <a-tag :color="loglevelColor[text] || 'blue'">{{ text.toUpperCase() }}</a-tag>
                </template>
                <template v-if="column.key === 'message'">
                    <a-tag v-if="record.actionDriver">{{record.actionDriver}}</a-tag>
                    <LogId :id="record.logId" v-if="record.type === 'action'" :color="pickLogIdColor(record.logId)"></LogId>
                    {{text}}
                </template>
            </template>
        </a-table>
    </a-drawer>
</template>

<style scoped>

</style>
