<script setup>
import {onMounted, ref} from "vue";
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

const router = useRouter()
const {httpGet} = useAxios()

const hooks = ref([])

const gotoHook = () => {
    router.push({name: "hooks"})
}
const fetchWebhooks = () => {
    httpGet("/webhook").exec().then(({payload}) => {
        hooks.value = payload
    }).catch(e => {

    })
}

const formatHookLink = (hook) => {
    return `${location.protocol}//${location.host}/hook/${hook.id}`
}

const copyHookLinkToClipboard = (hook) => {
    hook.copied = true
    navigator.clipboard.writeText(formatHookLink(hook)).then(() => {

    },() => {
    })
    setTimeout(() => {
        hook.copied = false
    }, 3000)
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
                <a-space direction="vertical" @click="copyHookLinkToClipboard(hook)">
                    <a-typography-text strong>
                        {{ hook.name }}
                        <a-button size="small" type="text" @click="copyHookLinkToClipboard(hook)">
                            <template #icon>
                                <CopyOutlined v-show="!hook.copied" />
                                <CheckOutlined color="success" v-show="hook.copied" />
                            </template>
                        </a-button>
                    </a-typography-text>
                    <a-typography-text>{{ hook.description }}</a-typography-text>
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
                    <span><a-tag color="blue">Run count:</a-tag> {{ hook.runCount }}</span>
                    <span v-if="hook.lastRunAt"><a-tag color="cyan">Last run:</a-tag>{{ hook.lastRunAt }}</span>
                </a-space>
            </a-col>
            <a-col :span="2">
                <a-space class="actions">
                    <a-button size="small" type="text">
                        <template #icon>
                            <InfoCircleOutlined/>
                        </template>
                    </a-button>
                    <a-button size="small" type="text">
                        <template #icon>
                            <EditOutlined/>
                        </template>
                    </a-button>
                    <DuplicateIt :id="hook.id" @duplicated="fetchWebhooks"></DuplicateIt>
                    <DeleteIt :id="hook.id" @deleted="fetchWebhooks"></DeleteIt>
                </a-space>
            </a-col>
        </a-row>
    </div>
</template>

<style scoped>

</style>
