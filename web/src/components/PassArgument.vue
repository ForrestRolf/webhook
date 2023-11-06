<script setup>
import {DeleteOutlined} from "@ant-design/icons-vue";
import {computed} from "vue";
import argumentSources from "../support/argument-source.js";

const emit = defineEmits(["update:argument", "remove"])
const props = defineProps({
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    },
    argument: {
        type: Object,
        default() {
            return {}
        }
    }
})

const argument = computed({
    get() {
        return props.argument
    },
    set(v) {
        emit("update:argument", v)
    }
})

const handleRemove = () => {
    emit("remove")
}
</script>

<template>
    <a-space>
        <a-select v-model:value="argument.source" :disabled="props.disabled">
            <a-select-option :value="src.source" v-for="src in argumentSources">{{ src.label }}</a-select-option>
        </a-select>
        <a-input v-model:value="argument.name" :disabled="props.disabled" placeholder="pusher.email"></a-input>
        <div>Set to env:</div>
        <a-input v-model:value="argument.envname" :disabled="props.disabled" placeholder="Optional"></a-input>

        <a-button type="text" danger @click="handleRemove" v-show="!props.disabled">
            <template #icon>
                <DeleteOutlined/>
            </template>
        </a-button>
    </a-space>
</template>

<style scoped>

</style>