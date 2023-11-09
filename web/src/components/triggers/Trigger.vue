<script setup>
import {computed} from "vue";
import {DeleteOutlined} from "@ant-design/icons-vue"
import argumentSources from "../../support/argument-source.js";

const props = defineProps({
    trigger: {
        type: Object,
        default() {
            return {}
        }
    },
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    }
})
const emit = defineEmits(["update:trigger", "remove"])

const trigger = computed({
    get() {
        return props.trigger
    },
    set(t) {
        emit("update:trigger", t)
    },
})
const handleRemove = () => {
    emit("remove")
}
</script>

<template>
    <a-space>
        <div>IF</div>
        <a-select v-model:value="trigger.match.parameter.source" :disabled="props.disabled">
            <a-select-option :value="src.name" v-for="src in argumentSources">{{ src.label }}</a-select-option>
        </a-select>
        <a-input v-model:value="trigger.match.parameter.name" :disabled="props.disabled"></a-input>

        <a-checkbox v-model:checked="trigger.match.not" :disabled="props.disabled">Not</a-checkbox>
        <a-select v-model:value="trigger.match.type" :disabled="props.disabled">
            <a-select-option value="value">Equal value</a-select-option>
            <a-select-option value="regex">Regular match</a-select-option>
        </a-select>

        <a-input v-model:value="trigger.match.value" v-show="trigger.match.type === 'value'" :disabled="props.disabled"></a-input>
        <a-input v-model:value="trigger.match.regex" v-show="trigger.match.type === 'regex'" :disabled="props.disabled"></a-input>

        <a-button type="text" danger @click="handleRemove" v-show="!props.disabled">
            <template #icon>
                <DeleteOutlined/>
            </template>
        </a-button>
    </a-space>
</template>

<style scoped>

</style>