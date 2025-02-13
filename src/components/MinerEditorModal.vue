<template>
    <div class="modal fade" id="minerEditor" tabindex="-1" role="dialog" aria-labelledby="minerEditorLabel"
         aria-hidden="true">
        <div class="modal-dialog modal-lg" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="minerEditorLabel">Miner Editor</h5>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close" @click="close">
                        <span aria-hidden="true">×</span>
                    </button>
                </div>
                <div class="modal-body">
                    <div class="form-group" v-if="selectedClient" :data-client-id="selectedClient.client_status.client_id">
                        <label for="config">Config for: {{ selectedClient.client_status.client_id }}</label>
                        <textarea class="form-control" rows="20" id="config" v-model="clientConfig"></textarea>
                    </div>
                    <div v-else>
                        Kein Client ausgewählt.
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-success" @click="saveConfig" :disabled="!selectedClient">
                        Speichern
                    </button>
                    <button type="button" class="btn btn-secondary" data-dismiss="modal" @click="close">Abbrechen
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';
import $ from 'jquery';

export default {
    props: {
        selectedClient: {
            type: Object,
            default: null,
        },
    },
    emits: ['close', 'configSaved', 'setStatusBar'],
    setup(props, { emit }) {
        const clientConfig = ref('');

        watch(() => props.selectedClient, async (newClient) => {
            if (newClient) {
                await fetchConfig();
            }
        });

        onMounted(() => {
            $('#minerEditor').on('hidden.bs.modal', () => {
                emit('close');
            });
        });


        const fetchConfig = async () => {
            try {
                const response = await axios.get(`/admin/getClientConfig?clientId=${props.selectedClient.client_status.client_id}`);
                clientConfig.value = JSON.stringify(response.data, null, 2);
            } catch (error) {
                console.error("Error fetching config:", error);
                emit('setStatusBar', `Fehler beim Abrufen der Konfiguration: ${error}`, 'danger');
            }
        };

        const saveConfig = async () => {
            try {
                await axios.post(`/admin/setClientConfig?clientId=${props.selectedClient.client_status.client_id}`, clientConfig.value);
                emit('configSaved');
                emit('setStatusBar', `Konfiguration erfolgreich gespeichert für Client: ${props.selectedClient.client_status.client_id}`, 'success');
                close();
            } catch (error) {
                console.error("Error saving config:", error);
                emit('setStatusBar', `Fehler beim Speichern der Konfiguration: ${error}`, 'danger');
            }
        };

        const close = () => {
            emit('close');
            $('#minerEditor').modal('hide');
        };

        return {
            clientConfig,
            saveConfig,
            close,
            fetchConfig
        };
    }
};
</script>