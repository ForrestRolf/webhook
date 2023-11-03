<script setup>
import {FileProtectOutlined} from "@ant-design/icons-vue"
import {computed} from "vue";

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
            </a-space>
        </a-col>
    </a-row>
</template>

<style lang="less">

</style>
