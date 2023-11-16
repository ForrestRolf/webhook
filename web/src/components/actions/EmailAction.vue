<script setup>
import {computed, onMounted, reactive, ref, watch} from "vue";
import useAxios from "../../support/axios.js";
import {useRoute} from "vue-router";
import _ from "lodash"
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

const route = useRoute()
const {httpGet} = useAxios()
const attributes = computed({
    get() {
        return props.attributes
    },
    set(v) {
        emit("update:attributes", v)
    }
})
const profiles = ref([])
const lang = ref("html")
const templatePickerRef = ref()

const fetchProfiles = () => {
    httpGet('/smtp/profile').exec().then((data) => {
        profiles.value = data.payload
    }).catch(e => {

    })
}

const handleCodeChange = (code) => {
    attributes.value.body = code
    emit("update:attributes", attributes.value)
}
const openCodeEditor = () => {
    props.handleCodeEditor({
        lang: lang.value,
        code: attributes.value.body || "",
        onSave: handleCodeChange
    })
}

const openTemplatePicker = () => {
    templatePickerRef.value.open()
}
const handleTemplateSelected = (template) => {
    attributes.value.body = template.content
    emit("update:attributes", attributes.value)
}

onMounted(() => {
    fetchProfiles()
})
</script>

<template>
    <a-row class="action-item">
        <div class="type-label">Email</div>
        <a-col :span="24">
            <a-form
                :label-col="{ span: 4 }"
                :wrapper-col="{ span: 20 }"
                :disabled="props.disabled"
            >

                <a-form-item label="SMTP profile">
                    <a-select v-model:value="attributes.profileId">
                        <a-select-option v-for="profile in profiles" :value="profile.id">{{profile.name}}</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item label="To">
                    <a-select
                        v-model:value="attributes.to"
                        mode="tags"
                        placeholder=""
                    ></a-select>
                </a-form-item>
                <a-form-item label="Cc">
                    <a-select
                        v-model:value="attributes.cc"
                        mode="tags"
                        placeholder=""
                    ></a-select>
                </a-form-item>
                <a-form-item label="Subject">
                    <a-input v-model:value="attributes.subject"></a-input>
                </a-form-item>
                <a-form-item label="Message">
                    <a-textarea :rows="3" v-model:value="attributes.body"></a-textarea>
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

        <TemplatePicker ref="templatePickerRef" lang="html" @selected="handleTemplateSelected"></TemplatePicker>
    </a-row>
</template>

<style scoped>

</style>
