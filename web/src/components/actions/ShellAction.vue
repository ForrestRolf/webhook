<script setup>
import {computed, ref} from "vue";
import {CodeOutlined, FileSearchOutlined} from "@ant-design/icons-vue"
import TemplatePicker from "../TemplatePicker.vue";

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
    arguments: {
        type: Array,
        default() {
            return []
        }
    },
    handleCodeEditor: {
        type: Function,
        default() {
            return () => {
            }
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
const templatePickerRef = ref()

const handleCodeChange = (code) => {
    attributes.value.scripts = code
    emit("update:attributes", attributes.value)
}

const openCodeEditor = () => {
    props.handleCodeEditor({
        lang: "shell",
        code: attributes.value.scripts || "#!/bin/bash",
        onSave: handleCodeChange
    })
}

const openTemplatePicker = () => {
    templatePickerRef.value.open()
}
const handleTemplateSelected = (template) => {
    attributes.value.scripts = template.content
    emit("update:attributes", attributes.value)
}
</script>

<template>
    <a-row class="action-item">
        <div class="type-label">Shell</div>
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="Working directory">
                    <a-input v-model:value="attributes.workingDirectory"></a-input>
                </a-form-item>
                <a-form-item label="Scripts">
                    <a-textarea :rows="8" v-model:value="attributes.scripts"></a-textarea>
                    <a-button size="small" type="text" @click="openCodeEditor">
                        <template #icon>
                            <CodeOutlined/>
                        </template>
                        Open in code editor
                    </a-button>
                    <a-divider type="vertical"></a-divider>
                    <a-button size="small" type="text" @click="openTemplatePicker">
                        <template #icon>
                            <FileSearchOutlined />
                        </template>
                        Start with a template
                    </a-button>
                </a-form-item>
            </a-form>
        </a-col>

        <TemplatePicker ref="templatePickerRef" lang="shell" @selected="handleTemplateSelected"></TemplatePicker>
    </a-row>
</template>

<style scoped>

</style>
