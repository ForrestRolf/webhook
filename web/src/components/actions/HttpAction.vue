<script setup>
import {computed} from "vue";

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
</script>

<template>
    <a-row class="action-item">
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="URL">
                    <a-input-group compact>
                        <a-select v-model:value="attributes.method">
                            <a-select-option value="GET">GET</a-select-option>
                            <a-select-option value="POST">POST</a-select-option>
                            <a-select-option value="PUT">PUT</a-select-option>
                        </a-select>
                        <a-input v-model:value="attributes.url" style="width: 50%" />
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
    </a-row>
</template>

<style scoped>

</style>
