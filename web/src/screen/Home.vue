<script setup>
import {computed, nextTick, onMounted, ref} from "vue";
import {
    PlusOutlined,
    EditOutlined,
    CopyOutlined,
    InfoCircleOutlined,
    CheckOutlined
} from '@ant-design/icons-vue';
import {useRouter} from "vue-router";
import useAxios from "../support/axios.js";
import TriggerGroupPreview from "../components/triggers/TriggerGroupPreview.vue";
import ActionPreview from "../components/actions/ActionPreview.vue";
import DuplicateIt from "../components/DuplicateIt.vue";
import DeleteIt from "../components/DeleteIt.vue";
import EnableOrDisableIt from "../components/EnableOrDisableIt.vue";
import WebhookLogs from "../components/WebhookLogs.vue";
import useMessage from "../support/message.js";

const router = useRouter()
const {httpGet} = useAxios()
const {successMessage} = useMessage()

const hooks = ref([])
const selectedHook = ref()
const logs = ref()
const loading = ref(true)
const keyword = ref(null)

const gotoHook = () => {
    router.push({name: "hooks"})
}
const fetchWebhooks = () => {
    loading.value = true
    httpGet("/webhook").exec().then(({payload}) => {
        hooks.value = payload
    }).catch(e => {

    }).finally(() => {
        loading.value = false
    })
}

const filteredHooks = computed(() => {
    if (!keyword.value) {
        return hooks.value
    }
    return hooks.value.filter(h => {
        return h.name.indexOf(keyword.value) > -1 || h.description.indexOf(keyword.value) > -1
    })
})

const formatHookLink = (hook) => {
    return `${location.protocol}//${location.host}/hook/${hook.id}`
}

const copyToClipboard = (msg) => {
    navigator.clipboard.writeText(msg).then(() => {
        successMessage("Copied").show()
    }, () => {
    })
}
const copyHookLinkToClipboard = (hook) => {
    hook.copied = true
    copyToClipboard(formatHookLink(hook))
    setTimeout(() => {
        hook.copied = false
    }, 3000)
}
const copyAuthToken = (hook) => {
    hook.authToken && copyToClipboard(`hook ${hook.authToken}`)
}

const handleSelectHook = (hook) => {
    selectedHook.value = hook.id
    nextTick(() => {
        logs.value.show()
    })
}

const gotoEdit = (hook) => {
    router.push({name: "hooks", query: {id: hook.id}})
}

onMounted(() => {
    setTimeout(fetchWebhooks, 100)
})
</script>

<template>
    <div class="home">
        <a-row class="new-hook">
            <a-col :span="12">
                <a-input-search
                    class="webhook-search"
                    v-model:value="keyword"
                    placeholder="Enter keywords to start searching"
                />
            </a-col>
            <a-col :span="12" class="txt-rgt">
                <a-button type="primary" @click="gotoHook">
                    <template #icon>
                        <PlusOutlined/>
                    </template>
                    New Hook
                </a-button>
            </a-col>
        </a-row>
        <a-row v-show="loading">
            <a-col :span="24">
                <a-skeleton active/>
            </a-col>
        </a-row>
        <a-row :gutter="12" class="hooks" v-for="hook in filteredHooks" align="middle">
            <a-col :span="4">
                <a-badge status="success" v-if="hook.enabled"/>
                <a-badge status="error" v-if="!hook.enabled"/>
                <a-space direction="vertical">
                    <a-typography-text strong>
                        {{ hook.name }}
                        <a-button size="small" type="text" @click="copyHookLinkToClipboard(hook)">
                            <template #icon>
                                <CopyOutlined v-show="!hook.copied"/>
                                <CheckOutlined color="success" v-show="hook.copied"/>
                            </template>
                        </a-button>
                    </a-typography-text>
                    <a-typography-text>{{ hook.description }}</a-typography-text>
                    <a-tooltip>
                        <template #title>Click to copy token</template>
                        <a-tag v-show="hook.authToken" color="pink" class="copyable" @click="copyAuthToken(hook)">
                            Authorization=hook ******
                        </a-tag>
                    </a-tooltip>
                </a-space>
            </a-col>
            <a-col :span="8">
                <TriggerGroupPreview :triggers="hook.triggers"></TriggerGroupPreview>
            </a-col>
            <a-col :span="5">
                <ActionPreview :actions="hook.actions"></ActionPreview>
            </a-col>
            <a-col :span="1">
                <EnableOrDisableIt :id="hook.id" :enabled="hook.enabled" @changed="fetchWebhooks"></EnableOrDisableIt>
            </a-col>
            <a-col :span="4">
                <a-space direction="vertical">
                    <span><a-tag color="green">Run count:</a-tag> {{ hook.runCount }} / {{ hook.callCount }}</span>
                    <span v-if="hook.lastRunAt"><a-tag color="cyan">Last run:</a-tag>{{ hook.lastRunAt }}</span>
                </a-space>
            </a-col>
            <a-col :span="2">
                <a-space class="actions">
                    <a-button size="small" type="text" @click="handleSelectHook(hook)">
                        <template #icon>
                            <InfoCircleOutlined/>
                        </template>
                    </a-button>
                    <a-button size="small" type="text" @click="gotoEdit(hook)">
                        <template #icon>
                            <EditOutlined/>
                        </template>
                    </a-button>
                    <DuplicateIt :id="hook.id" @duplicated="fetchWebhooks"></DuplicateIt>
                    <DeleteIt :id="hook.id" @deleted="fetchWebhooks"></DeleteIt>
                </a-space>
            </a-col>
        </a-row>

        <WebhookLogs ref="logs" :webhook-id="selectedHook"></WebhookLogs>
    </div>
</template>

<style scoped>

</style>
