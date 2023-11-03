<script setup>
import {onMounted, reactive, ref} from "vue";
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    CopyOutlined,
    PlayCircleOutlined,
    PauseCircleOutlined,
    InfoCircleOutlined
} from '@ant-design/icons-vue';
import {useRouter} from "vue-router";
import useAxios from "../support/axios.js";
import TriggerGroupPreview from "../components/triggers/TriggerGroupPreview.vue";
import ActionPreview from "../components/actions/ActionPreview.vue";

const router = useRouter()
const {httpGet, httpPost, httpPut} = useAxios()

const hooks = ref([])
const loading = reactive({
    enable: false,
    disable: false,
    duplicate: false
})

const gotoHook = () => {
    router.push({name: "hooks"})
}
const fetchWebhooks = () => {
    httpGet("/webhook").exec().then(({payload}) => {
        hooks.value = payload
    }).catch(e => {

    })
}
const handleEnableHook = (id) => {
    loading.enable = true
    httpPut(`/webhook/${id}/enable`).exec().then(({payload}) => {
        fetchWebhooks()
    }).catch(e => {

    }).finally(() => {
        loading.enable = false
    })
}
const handleDisableHook = (id) => {
    loading.disable = true
    httpPut(`/webhook/${id}/disable`).exec().then(({payload}) => {
        fetchWebhooks()
    }).catch(e => {

    }).finally(() => {
        loading.disable = false
    })
}

onMounted(() => {
    fetchWebhooks()
})
</script>

<template>
    <div class="home">
        <a-row class="new-hook">
            <a-col :span="24">
                <a-button type="primary" @click="gotoHook">
                    <template #icon>
                        <PlusOutlined/>
                    </template>
                    New Hook
                </a-button>
            </a-col>
        </a-row>
        <a-row :gutter="12" class="hooks" v-for="hook in hooks" align="middle">
            <a-col :span="4">
                <a-badge status="success" v-if="hook.enabled"/>
                <a-badge status="error" v-if="!hook.enabled"/>
                <a-space direction="vertical">
                    <span>{{ hook.name }}</span>
                    <a-typography-text type="secondary">{{ hook.description }}</a-typography-text>
                </a-space>
            </a-col>
            <a-col :span="8">
                <TriggerGroupPreview :triggers="hook.triggers"></TriggerGroupPreview>
            </a-col>
            <a-col :span="5">
                <ActionPreview :actions="hook.actions"></ActionPreview>
            </a-col>
            <a-col :span="1">
                <a-button type="text" size="large" v-if="hook.enabled" :loading="loading.disable"
                          @click="handleDisableHook(hook.id)">
                    <template #icon>
                        <PlayCircleOutlined color="primary"/>
                    </template>
                </a-button>
                <a-button type="text" danger size="large" v-if="!hook.enabled" :loading="loading.enable"
                          @click="handleEnableHook(hook.id)">
                    <template #icon>
                        <PauseCircleOutlined/>
                    </template>
                </a-button>
            </a-col>
            <a-col :span="4">
                <a-space direction="vertical">
                    <span><a-tag color="blue">Run count:</a-tag> {{ hook.runCount }}</span>
                    <span><a-tag color="cyan">Last run:</a-tag>{{ hook.lastRunAt }}</span>
                </a-space>
            </a-col>
            <a-col :span="2">
                <a-space class="actions">
                    <a-button size="small" type="text">
                        <template #icon>
                            <InfoCircleOutlined />
                        </template>
                    </a-button>
                    <a-button size="small" type="text">
                        <template #icon>
                            <EditOutlined/>
                        </template>
                    </a-button>
                    <a-button size="small" type="text">
                        <template #icon>
                            <CopyOutlined color="blue"/>
                        </template>
                    </a-button>
                    <a-button size="small" danger type="text">
                        <template #icon>
                            <DeleteOutlined/>
                        </template>
                    </a-button>
                </a-space>
            </a-col>
        </a-row>
    </div>
</template>

<style scoped>

</style>
