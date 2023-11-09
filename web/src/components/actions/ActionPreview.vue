<script setup>
import {FileProtectOutlined, CopyOutlined} from "@ant-design/icons-vue"
import {computed} from "vue";
import useClipboard from "../../support/clipboard.js";


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
</script>

<template>
    <a-row :gutter="[0, 5]">
        <a-col :span="24" v-for="(action, i) in actions">
            <a-space class="action-preview">
                <template v-if="action.driver === 'shell'">
                    <span>{{ i + 1 }}.</span>
                    <span>Exec shell in directory: {{ action.attributes?.workingDirectory }}</span>
                    <a-tooltip>
                        <template #title>{{ action.attributes?.scripts }}</template>
                        <a-tag color="orange">
                            <FileProtectOutlined/>
                        </a-tag>
                    </a-tooltip>
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
            </a-space>
        </a-col>
    </a-row>
</template>

<style lang="less">

</style>
