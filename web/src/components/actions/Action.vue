<script setup>
import {computed, ref} from "vue";
import {
    ArrowDownOutlined,
    DeleteOutlined,
    NodeExpandOutlined,
    BranchesOutlined,
    CodeOutlined,
    SlackOutlined,
    MailOutlined,
    WechatOutlined,
    DingdingOutlined,
    CommentOutlined
} from "@ant-design/icons-vue"
import ShellAction from "./ShellAction.vue";
import HttpAction from "./HttpAction.vue";
import CodeEditor from "../CodeEditor.vue";
import Dispatcher from "./Dispatcher.vue";
import EmailAction from "./EmailAction.vue";
import SlackAction from "./SlackAction.vue";
import SmsAction from "./SmsAction.vue";

const emit = defineEmits(["update:actions"])
const props = defineProps({
    actions: {
        type: Object,
        default() {
            return []
        }
    },
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    },
    arguments: {
        type: Array,
        default() {
            return []
        }
    }
})

const actions = computed({
    get() {
        return props.actions
    },
    set(v) {
        emit("update:actions", v)
    }
})

const components = {
    "shell": ShellAction,
    "http": HttpAction,
    "dispatcher": Dispatcher,
    "email": EmailAction,
    "slack": SlackAction,
    "sms-twilio": SmsAction,
    "sms-plivo": SmsAction,
    "sms-burst": SmsAction,
    "sms-sns": SmsAction,
}
const handleRemove = (idx) => {
    actions.value.splice(idx, 1)
}
const addShellAction = () => {
    actions.value.push({
        "driver": "shell",
        "attributes": {
            "workingDirectory": "/tmp"
        }
    })
    emit("update:actions", actions.value)
}
const addHttpAction = () => {
    actions.value.push({
        "driver": "http",
        "attributes": {
            "method": "POST",
            "contentType": "application/json",
            "timeout": 30,
            "saveResponse": true,
        }
    })
    emit("update:actions", actions.value)
}
const addDispatcherAction = () => {
    actions.value.push({
        "driver": "dispatcher",
        "attributes": {
            "if": {},
            "url": "",
            "method": "POST",
            "webhookId": "",
            "webhookName": "",
            "compare": "eq",
        }
    })
    emit("update:actions", actions.value)
}

const addEmailAction = () => {
    actions.value.push({
        "driver": "email",
        "attributes": {
            profileId: "",
            to: [],
            cc: [],
            subject: "",
            body: "",
        }
    })
    emit("update:actions", actions.value)
}

const addSlackAction = () => {
    actions.value.push({
        "driver": "slack",
        "attributes": {
            webhookUrl: "",
            message: "",
        }
    })
    emit("update:actions", actions.value)
}

const handleSmsAction = (act) => {
    actions.value.push({
        "driver": act.key,
        "attributes": {
            profileId: "",
            provider: act.key,
            to: "",
            content: "",
        }
    })
    emit("update:actions", actions.value)
}

const codeEditor = ref()
const onCodeSave = ref()
const handleCodeEditor = ({code, lang, onSave}) => {
    onCodeSave.value = onSave
    codeEditor.value.open(lang, code)
}
</script>

<template>
    <a-row :gutter="[12, 12]" class="actions">
        <a-col :span="24" v-for="(action, i) in actions">
            <a-button class="float-rgt" type="text" danger @click="handleRemove(i)" v-show="!props.disabled">
                <template #icon>
                    <DeleteOutlined/>
                </template>
            </a-button>
            <component
                :is="components[action.driver]"
                :disabled="props.disabled"
                v-model:attributes="actions[i].attributes"
                :arguments="props.arguments"
                :handle-code-editor="handleCodeEditor">
            </component>
            <a-divider v-show="i < actions.length - 1">
                <ArrowDownOutlined/>
            </a-divider>
        </a-col>
        <a-col :span="24" class="txt-center" v-show="!props.disabled">
            <a-space>
                <span>Add</span>
                <a-tooltip title="Shell script">
                    <a-button type="text" @click="addShellAction"><template #icon><CodeOutlined/></template></a-button>
                </a-tooltip>
                <a-tooltip title="HTTP request">
                    <a-button type="text" @click="addHttpAction"><template #icon><NodeExpandOutlined/></template></a-button>
                </a-tooltip>
                <a-tooltip title="Webhook dispatcher">
                    <a-button type="text" @click="addDispatcherAction"><template #icon><BranchesOutlined/></template></a-button>
                </a-tooltip>
                <a-tooltip title="Slack" @click="addSlackAction">
                    <a-button type="text"><template #icon><SlackOutlined /></template></a-button>
                </a-tooltip>
                <a-tooltip title="Email" @click="addEmailAction">
                    <a-button type="text"><template #icon><MailOutlined /></template></a-button>
                </a-tooltip>
                <a-tooltip title="SMS">
                    <a-dropdown>
                        <template #overlay>
                            <a-menu @click="handleSmsAction">
                                <a-menu-item key="sms-twilio">
                                    Twilio
                                </a-menu-item>
                                <a-menu-item key="sms-burst">
                                    Burst SMS
                                </a-menu-item>
                                <a-menu-item key="sms-plivo">
                                    Plivo
                                </a-menu-item>
                                <a-menu-item key="sms-sns">
                                    SNS
                                </a-menu-item>
                            </a-menu>
                        </template>
                        <CommentOutlined />
                    </a-dropdown>
                </a-tooltip>
                <a-tooltip title="Wechat (coming soon)">
                    <a-button type="text"><template #icon><WechatOutlined /></template></a-button>
                </a-tooltip>
                <a-tooltip title="Dingding (coming soon)">
                    <a-button type="text"><template #icon><DingdingOutlined /></template></a-button>
                </a-tooltip>
            </a-space>
        </a-col>

        <CodeEditor ref="codeEditor" :on-save="onCodeSave"></CodeEditor>
    </a-row>
</template>

<style scoped>

</style>
