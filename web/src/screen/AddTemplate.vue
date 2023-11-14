<script setup>
import {SaveOutlined, CloseOutlined} from "@ant-design/icons-vue"
import CodeEditor from "../components/CodeEditor.vue";
import {computed, reactive, ref, watch} from "vue";
import useAxios from "../support/axios.js";
import useMessage from "../support/message.js";

const {httpPost, httpPut} = useAxios()
const {successMessage, errorMessage} = useMessage()
const emit = defineEmits(["created"])

const codeEditor = ref()
const formState = reactive({
    title: "",
    description: "",
    language: "shell",
    content: "",
})
const saving = ref(false)
const template = ref()
const isEditMode = ref(false)
const language = computed(() => {
    return formState.language
})

const open = (temp) => {
    isEditMode.value = !!temp
    template.value = temp

    if (!isEditMode.value) {
        temp = {
            title: "",
            description: "",
            language: "shell",
            content: "",
        }
    }
    formState.title = temp.title
    formState.language = temp.language
    formState.description = temp.description
    codeEditor.value.open(temp.language, temp.content)
}
const close = () => {
    codeEditor.value.close()
}
const handleSave = () => {
    formState.content = codeEditor.value.getCode()

    if (formState.title === "" || formState.content === "") {
        errorMessage("Could not save empty template").show()
        return
    }

    let client
    if (isEditMode.value) {
        client = httpPut(`/template/${template.value.id}`)
    } else {
        client = httpPost('/template')
    }
    saving.value = true
    client.withBody(formState).exec().then(() => {
        successMessage(`Template ${isEditMode ? "updated" : "created"}`).show()
        close()
        emit("created")
    }).catch(e => {
        errorMessage(`Could not ${isEditMode ? "update" : "save"} template`, e).show()
    }).finally(() => {
        saving.value = false
    })
}

watch(language, () => {
    codeEditor.value.setLanguage(language.value)
})

defineExpose({open, close})
</script>

<template>
    <CodeEditor ref="codeEditor" :on-save="handleSave" hide-default-actions title="Template Editor">
        <a-form class="add-template-form">
            <a-form-item label="Title">
                <a-input v-model:value="formState.title"></a-input>
            </a-form-item>
            <a-form-item label="Description">
                <a-textarea :rows="3" v-model:value="formState.description"></a-textarea>
            </a-form-item>
            <a-form-item label="Language">
                <a-select v-model:value="formState.language">
                    <a-select-option value="json">JSON</a-select-option>
                    <a-select-option value="shell">Shell</a-select-option>
                    <a-select-option value="plaintext">Text</a-select-option>
                    <a-select-option value="xml">XML</a-select-option>
                </a-select>
            </a-form-item>
            <a-form-item label="Template"></a-form-item>
        </a-form>

        <template #actions>
            <a-space>
                <a-button @click="close">
                    <template #icon>
                        <CloseOutlined/>
                    </template>
                    Close
                </a-button>
                <a-button type="primary" @click="handleSave" :loading="saving">
                    <template #icon>
                        <SaveOutlined/>
                    </template>
                    {{ isEditMode ? "Update" : "Save" }}
                </a-button>
            </a-space>
        </template>
    </CodeEditor>
</template>

<style lang="less">
.add-template-form {
    .ant-form-item-label {
        width: 100px;
        text-align: left;
    }

    .ant-select {
        width: 200px;
    }

    input.ant-input {
        width: 30%;
    }
}
</style>
