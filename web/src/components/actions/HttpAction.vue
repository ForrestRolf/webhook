<script setup>
import {computed, ref} from "vue";
import {CodeOutlined, FileSearchOutlined} from "@ant-design/icons-vue";
import TemplatePicker from "../TemplatePicker.vue";

const props = defineProps({
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    },
    attributes: {
        type: Object,
        default() {
            return {}
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
            return () => {}
        }
    }
})

const emit = defineEmits(["update:attributes"])
const templatePickerRef = ref()

const attributes = computed({
    get() {
        return props.attributes
    },
    set(v) {
        emit("update:attributes", v)
    }
})
const handleCodeChange = (code) => {
    attributes.value.payload = code
    emit("update:attributes", attributes.value)
}
const lang = computed(() => {
    let mime2language = {
        "text/plain": "plaintext",
        "application/json": "json",
        "application/xml": "xml",
    }
    return mime2language[attributes.value.contentType] ? mime2language[attributes.value.contentType] : "plaintext"
})

const openCodeEditor = () => {
    props.handleCodeEditor({
        lang: lang.value,
        code: attributes.value.payload || defaultCode[lang.value],
        onSave: handleCodeChange
    })
}

const openTemplatePicker = () => {
    templatePickerRef.value.open()
}
const handleTemplateSelected = (template) => {
    attributes.value.payload = template.content
    emit("update:attributes", attributes.value)
}
</script>

<template>
    <a-row class="action-item">
        <div class="type-label">HTTP</div>
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="URL">
                    <a-input-group compact>
                        <a-select v-model:value="attributes.method">
                            <a-select-option value="GET">GET</a-select-option>
                            <a-select-option value="POST">POST</a-select-option>
                            <a-select-option value="PUT">PUT</a-select-option>
                        </a-select>
                        <a-input v-model:value="attributes.url" style="width: 50%"/>
                    </a-input-group>
                </a-form-item>
                <a-form-item label="Content type">
                    <a-select v-model:value="attributes.contentType">
                        <a-select-option value="text/plain">Text</a-select-option>
                        <a-select-option value="application/json">JSON</a-select-option>
                        <a-select-option value="application/xml">XML</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item label="Payload">
                    <a-textarea v-model:value="attributes.payload" :rows="4"></a-textarea>
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
                <a-form-item label="Auth header">
                    <a-input v-model:value="attributes.authToken"></a-input>
                </a-form-item>
                <a-form-item label="Timeout">
                    <a-space>
                        <a-input-number v-model:value="attributes.timeout"></a-input-number>
                        <a-checkbox v-model:checked="attributes.saveResponse">Save response</a-checkbox>
                    </a-space>
                </a-form-item>
            </a-form>
        </a-col>

        <TemplatePicker ref="templatePickerRef" :lang="lang" @selected="handleTemplateSelected"></TemplatePicker>
    </a-row>
</template>

<style scoped>

</style>
