<script setup>
import {computed, onMounted, reactive, ref, watch} from "vue";
import useAxios from "../../support/axios.js";
import {useRoute} from "vue-router";
import _ from "lodash"

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
const currentWebhookId = computed(() => {
    return route.query.id
})

const ifArgument = reactive({
    name: "",
    value: ""
})
const webhooks = ref([])
const webhookId = ref()
const fetchWebhooks = () => {
    httpGet(`/webhook`).exec().then(({payload}) => {
        webhooks.value = payload
    }).catch(e => {

    }).finally(() => {

    })
}

watch(ifArgument, () => {
    if (ifArgument.name) {
        attributes.value.if[ifArgument.name] = ifArgument.value
    } else {
        attributes.value.if = {}
    }
    emit("update:attributes", attributes.value)
})
watch(attributes, () => {
    if (attributes.value.if) {
        let keys = Object.keys(attributes.value.if)
        ifArgument.name = keys[0] || ""
        ifArgument.value = attributes.value.if[keys[0]] || ""
    }
    webhookId.value = attributes.value.webhookId
}, {
    immediate: true
})
watch(webhookId, () => {
    if (!webhookId.value) return
    const url = import.meta.env.VITE_API_URL || location.href
    let u = new URL(url)

    const webhook = _.find(webhooks.value, {id: webhookId.value})
    attributes.value.webhookName = webhook.name
    attributes.value.webhookId = webhookId.value
    attributes.value.url = `${u.protocol}//${u.host}/hook/${webhookId.value}`
    emit("update:attributes", attributes.value)
})
onMounted(() => {
    fetchWebhooks()
})
</script>

<template>
    <a-row class="action-item">
        <div class="type-label">Dispatcher</div>
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="If">
                    <a-space>
                        <a-select v-model:value="ifArgument.name">
                            <a-select-option value="">Any</a-select-option>
                            <a-select-option v-for="arg in props.arguments" :value="arg.name">
                                {{ arg.name }}
                            </a-select-option>
                        </a-select>
                        <span v-if="ifArgument.name">=</span>
                        <a-input v-if="ifArgument.name" v-model:value="ifArgument.value"></a-input>
                    </a-space>
                </a-form-item>
                <a-form-item label="URL">
                    <a-input-group compact>
                        <a-select v-model:value="attributes.method">
                            <a-select-option value="GET">GET</a-select-option>
                            <a-select-option value="POST">POST</a-select-option>
                            <a-select-option value="PUT">PUT</a-select-option>
                        </a-select>
                        <a-select v-model:value="webhookId" class="webhook-select">
                            <a-select-option v-for="hook in webhooks" :value="hook.id"
                                             :disabled="currentWebhookId === hook.id">{{ hook.name }}
                            </a-select-option>
                        </a-select>
                    </a-input-group>
                </a-form-item>
            </a-form>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>