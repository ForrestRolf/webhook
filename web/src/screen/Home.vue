<script setup>
import {computed, nextTick, onMounted, ref, watch} from "vue";
import {
    PlusOutlined,
    EditOutlined,
    CopyOutlined,
    InfoCircleOutlined,
    CheckOutlined,
    CloudDownloadOutlined,
    CloudUploadOutlined,
    SyncOutlined,
    SortAscendingOutlined
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
import {download} from "../support/file.js";
import useClipboard from "../support/clipboard.js";
import CodeEditor from "../components/CodeEditor.vue";

const router = useRouter()
const {httpGet} = useAxios()
const {successMessage, errorMessage} = useMessage()
const {copyToClipboard} = useClipboard()

const hooks = ref([])
const selectedHook = ref()
const logs = ref()
const loading = ref(true)
const keyword = ref(null)
const codeEditor = ref()
const orderBy = ref("lastRun")

const gotoHook = () => {
    router.push({name: "hooks"})
}
const fetchWebhooks = () => {
    loading.value = true
    let query = {
        orderBy: orderBy.value
    }
    httpGet("/webhook").withQuery(query).exec().then(({payload}) => {
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

const fileList = ref([])
const uploadUrl = computed(() => {
    return `${import.meta.env.VITE_API_URL || "/api"}/import`
})

const formatHookLink = (hook) => {
    const h = new URL(import.meta.env.VITE_API_URL || location.href)
    return `${h.protocol}//${h.host}/hook/${hook.id}`
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
const handleExport = (hook) => {
    httpGet(`/webhook/${hook.id}`).exec().then(({payload}) => {
        download(JSON.stringify(payload), `${hook.id}.json`, "application/json")
    }).catch()
}
const onUploadFileChange = ({file}) => {
    switch (file.status) {
        case "done":
            if (file.response.meta.code === "OK") {
                successMessage("Import successfully").show()
                fetchWebhooks()
            }
            break;
        case "error":
            errorMessage("Unable to import this file. Please try again", file?.response?.meta?.message).show()
            break;
    }
}

const gotoEdit = (hook) => {
    router.push({name: "hooks", query: {id: hook.id}})
}
const onCodePreview = ({lang, code}) => {
    codeEditor.value.open(lang, code)
}

watch(orderBy, () => {
    fetchWebhooks()
})
onMounted(() => {
    setTimeout(fetchWebhooks, 100)
})
</script>

<template>
    <div class="home">
        <a-row class="new-hook">
            <a-col :span="12">
                <a-space>
                    <a-input-search
                        class="webhook-search"
                        v-model:value="keyword"
                        placeholder="Enter keywords to start searching"
                    />
                    <a-radio-group v-model:value="orderBy" button-style="solid">
                        <a-radio-button value="" disabled color="white">Order by:</a-radio-button>
                        <a-radio-button value="create">Create date</a-radio-button>
                        <a-radio-button value="update">Update date</a-radio-button>
                        <a-radio-button value="lastRun">Last run</a-radio-button>
                    </a-radio-group>
                </a-space>
            </a-col>
            <a-col :span="12" class="txt-rgt">
                <a-space>
                    <a-button @click="fetchWebhooks">
                        <template #icon>
                            <SyncOutlined />
                        </template>
                    </a-button>
                    <a-divider type="vertical" />
                    <a-upload
                        v-model:file-list="fileList"
                        name="file"
                        :action="uploadUrl"
                        :headers="{'X-Requested-With': null}"
                        @change="onUploadFileChange"
                    >
                        <a-button>
                            <template #icon>
                                <CloudUploadOutlined/>
                            </template>
                            Import
                        </a-button>
                    </a-upload>
                    <a-button type="primary" @click="gotoHook">
                        <template #icon>
                            <PlusOutlined/>
                        </template>
                        New Hook
                    </a-button>
                </a-space>
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
                    <a-typography-text v-if="hook.description">{{ hook.description }}</a-typography-text>
                    <a-tooltip v-if="hook.authToken">
                        <template #title>Click to copy token</template>
                        <a-tag v-show="hook.authToken" color="pink" class="copyable" @click="copyAuthToken(hook)">
                            Authorization=hook ******
                        </a-tag>
                    </a-tooltip>
                    <a-tag v-if="hook.debug" color="orange">Debug enabled</a-tag>
                </a-space>
            </a-col>
            <a-col :span="8">
                <TriggerGroupPreview :triggers="hook.triggers"></TriggerGroupPreview>
            </a-col>
            <a-col :span="5">
                <ActionPreview :actions="hook.actions" @code-preview="onCodePreview"></ActionPreview>
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
                    <a-button size="small" type="text" @click="handleExport(hook)">
                        <template #icon>
                            <CloudDownloadOutlined/>
                        </template>
                    </a-button>
                    <DuplicateIt :id="hook.id" @duplicated="fetchWebhooks"></DuplicateIt>
                    <DeleteIt :id="hook.id" @deleted="fetchWebhooks"></DeleteIt>
                </a-space>
            </a-col>
        </a-row>

        <WebhookLogs ref="logs" :webhook-id="selectedHook"></WebhookLogs>
        <CodeEditor ref="codeEditor" read-only></CodeEditor>
    </div>
</template>

<style>
.ant-upload-list {
    display: none;
}
</style>
