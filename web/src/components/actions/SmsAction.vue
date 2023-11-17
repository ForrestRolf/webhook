<script setup>
import {computed, onMounted, ref} from "vue";
import {useRoute} from "vue-router";
import useAxios from "../../support/axios.js";

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
const profilesByProvider = computed(() => {
    return profiles.value.filter(p => {
        return p.provider === props.attributes.provider
    })
})

const fetchProfiles = () => {
    httpGet('/sms/profile').exec().then((data) => {
        profiles.value = data.payload
    }).catch(e => {

    })
}

onMounted(() => {
    fetchProfiles()
})
</script>

<template>
    <a-row class="action-item">
        <div class="type-label">{{ $filters.smsProvider(attributes.provider) }}</div>
        <a-col :span="24">
            <a-form
                :label-col="{ span: 5 }"
                :wrapper-col="{ span: 20 }"
                :disabled="props.disabled"
            >
                <a-form-item label="Profile">
                    <a-select v-model:value="attributes.profileId">
                        <a-select-option v-for="profile in profilesByProvider" :value="profile.id">{{ profile.name }}
                        </a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item label="To">
                    <a-input v-model:value="attributes.to"></a-input>
                </a-form-item>
                <a-form-item label="Message">
                    <a-textarea :rows="6" v-model:value="attributes.content"></a-textarea>
                </a-form-item>
            </a-form>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
