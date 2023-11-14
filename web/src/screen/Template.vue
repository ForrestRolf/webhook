<script setup>
import {PlusOutlined} from "@ant-design/icons-vue"
import Gist from "../components/Gist.vue";
import {computed, onMounted, ref} from "vue";
import AddTemplate from "./AddTemplate.vue";
import useAxios from "../support/axios.js";
import useMessage from "../support/message.js";
import _ from "lodash"
import CodeEditor from "../components/CodeEditor.vue";

const {httpGet, httpDelete} = useAxios()
const {successMessage, errorMessage} = useMessage()

const templates = ref([])
const addTemplateRef = ref()
const codeEditorRef = ref()
const loading = ref(false)
const selectedTemplate = ref()
const keyword = ref()
const filteredTemplates = computed(() => {
    if (!keyword.value) {
        return templates.value
    }
    return templates.value.filter(t => {
        return t.title.indexOf(keyword.value) > -1 || t.description.indexOf(keyword.value) > -1
    })
})

const fetchTemplates = () => {
    loading.value = true
    httpGet("/template").exec().then((data) => {
        templates.value = data.payload
    }).finally(() => {
        loading.value = false
    })
}

const handleOpenTemplateEditor = () => {
    addTemplateRef.value.open()
}
const handleEdit = (template) => {
    selectedTemplate.value = template
    addTemplateRef.value.open(template)
}
const handleDelete = (template) => {
    let idx = _.findIndex(templates.value, {id: template.id})

    httpDelete(`/template/${template.id}`).exec().then(() => {
        templates.value.splice(idx, 1)
        successMessage(`Template deleted`, template.title).show()
    }).catch(e => {
        errorMessage("Could not delete template", e).show()
    }).finally(() => {

    })
}
const handleReadCode = (template) => {
    codeEditorRef.value.open(template.language, template.content)
}

onMounted(() => {
    fetchTemplates()
})
</script>

<template>
    <a-row :gutter="[0, 12]">
        <a-col :span="12">
            <a-input-search
                class="template-search"
                v-model:value="keyword"
                placeholder="Enter keywords to start searching"
            />
        </a-col>
        <a-col :span="12" class="txt-rgt">
            <a-button type="primary" @click="handleOpenTemplateEditor">
                <template #icon>
                    <PlusOutlined/>
                </template>
                Add new template
            </a-button>

            <AddTemplate ref="addTemplateRef" @created="fetchTemplates"></AddTemplate>
        </a-col>
        <a-col :span="24">
            <a-row :gutter="[12, 12]">
                <a-col :span="4" v-for="template in filteredTemplates">
                    <Gist :template="template" @edit="handleEdit" @delete="handleDelete" @read="handleReadCode"></Gist>
                </a-col>
                <a-col :span="24" v-if="filteredTemplates.length === 0">
                    <a-empty />
                </a-col>
            </a-row>
            <CodeEditor ref="codeEditorRef" read-only id="preview-editor"></CodeEditor>
        </a-col>
    </a-row>
</template>

<style lang="less">
    .template-search {
        width: 350px;
    }
</style>
