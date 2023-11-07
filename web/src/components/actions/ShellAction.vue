<script setup>
import {computed, ref} from "vue";
import {CodeOutlined} from "@ant-design/icons-vue"

const props = defineProps({
    attributes: {
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
    },
    handleCodeEditor: {
        type: Function,
        default() {
            return () => {}
        }
    }
})
const emit = defineEmits(["update:attributes"])

const attributes = computed({
    get() {
        return props.attributes
    },
    set(v) {
        emit("update:attributes", v)
    }
})

const handleCodeChange = (code) => {
    attributes.value.scripts = code
}

const openCodeEditor = () => {
    props.handleCodeEditor({
        lang: "shell",
        code: attributes.value.scripts || "#!/bin/bash",
        onSave: handleCodeChange
    })
}
</script>

<template>
    <a-row class="action-item">
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="Working directory">
                    <a-input v-model:value="attributes.workingDirectory"></a-input>
                </a-form-item>
                <a-form-item label="Scripts">
                    <a-textarea :rows="8" v-model:value="attributes.scripts"></a-textarea>
                    <a-button size="small" type="text" @click="openCodeEditor">
                        <template #icon>
                            <CodeOutlined />
                        </template>
                        Open in code editor
                    </a-button>
                </a-form-item>
            </a-form>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
