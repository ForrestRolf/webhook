<script setup>
import {computed} from "vue";
import {useRoute} from "vue-router";

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
})
const emit = defineEmits(["update:attributes"])

const route = useRoute()
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
        <div class="type-label">Slack</div>
        <a-col :span="24">
            <a-form
                :label-col="{ span: 4 }"
                :wrapper-col="{ span: 20 }"
                :disabled="props.disabled"
            >
                <a-form-item label="Webhook URL">
                    <a-input v-model:value="attributes.webhookUrl"></a-input>
                </a-form-item>
                <a-form-item label="Channel">
                    <a-input v-model:value="attributes.channel">
                        <template #prefix>#</template>
                    </a-input>
                </a-form-item>
                <a-form-item label="Message">
                    <a-textarea :rows="6" v-model:value="attributes.message"></a-textarea>
                </a-form-item>
            </a-form>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>