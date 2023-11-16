<script setup>
import {FileTextOutlined, CopyOutlined} from "@ant-design/icons-vue"
import {computed} from "vue";
import useClipboard from "../../support/clipboard.js";

const emit = defineEmits(["codePreview"])
const props = defineProps({
    actions: {
        type: Array,
        default() {
            return []
        }
    }
})
const actions = computed(() => {
    return props.actions
})
const {copyToClipboard} = useClipboard()
const copy = (url) => {
    copyToClipboard(url)
}
const handleOpenCodePreview = (lang, code) => {
    emit("codePreview", {
        lang,
        code
    })
}
</script>

<template>
    <a-row :gutter="[0, 5]">
        <a-col :span="24" v-for="(action, i) in actions">
            <a-space class="action-preview">
                <template v-if="action.driver === 'shell'">
                    <span>{{ i + 1 }}.</span>
                    <span>Exec shell in directory: {{ action.attributes?.workingDirectory }}</span>
                    <a-button @click="handleOpenCodePreview('shell', action.attributes?.scripts)" type="text">
                        <template #icon>
                            <FileTextOutlined />
                        </template>
                    </a-button>
                </template>
                <template v-if="action.driver === 'http'">
                    <span>{{ i + 1 }}.</span>
                    <a-tag color="blue">{{ action.attributes?.method }}</a-tag>
                    <span>{{ action.attributes?.url }}</span>
                </template>

                <template v-if="action.driver === 'dispatcher'">
                    <span>{{ i + 1 }}.</span>
                    <span v-if="Object.keys(action.attributes?.if).length > 0">If </span>
                    <a-tag v-if="Object.keys(action.attributes?.if).length > 0">
                        {{ Object.keys(action.attributes?.if).pop() }} = {{ Object.values(action.attributes?.if).pop() }}
                    </a-tag>
                    <span v-if="Object.keys(action.attributes?.if).length === 0">Always</span>
                    <span>re-send to</span>
                    <a-tag color="#2db7f5" class="copyable" @click="copy(action.attributes?.url)">
                        {{ action.attributes?.webhookName }}
                    </a-tag>
                </template>

                <template v-if="action.driver === 'email'">
                    <span>{{ i + 1 }}.</span>
                    <span class="flex">
                        <span>Send mail to: </span>
                        <span class="flex flex-col pad-lft">
                            <span v-for="t in action.attributes?.to">{{t}}</span>
                            <span v-if="action.attributes?.cc.length > 0" class="mgr-top">Cc: </span>
                            <span v-if="action.attributes?.cc.length > 0" v-for="t in action.attributes?.cc">{{t}}</span>
                        </span>
                    </span>
                </template>
            </a-space>
        </a-col>
    </a-row>
</template>

<style lang="less">
    .action-preview {
        align-items: start;
    }
</style>
