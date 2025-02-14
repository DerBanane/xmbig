import $ from 'jquery';
import numeral from 'numeral';

// Keine Importzeile für timeago

function clientStatus(data, TRESHOLD_IN_MS, currentServerTime, row) {
    var lastStatus = row.client_status.last_status_update * 1000;

    if (isOnline(lastStatus, TRESHOLD_IN_MS, currentServerTime)) {
        return data;
    } else {
        return "OFFLINE";
    }
}

function clientInfo(data, type, TRESHOLD_IN_MS, currentServerTime, row) {
    if (type !== 'sort') {
        var lastStatus = row.client_status.last_status_update * 1000;
        var online = isOnline(lastStatus, TRESHOLD_IN_MS, currentServerTime);

        var tooltip = "CPU: " + row.client_status.cpu_brand + " (" + row.client_status.cpu_sockets + ") [" + row.client_status.cpu_cores + " cores / " + row.client_status.cpu_threads + " threads]";
        tooltip += '\n';
        tooltip += "CPU Flags: " + (row.client_status.cpu_has_aes ? "AES-NI " : "");
        tooltip += (row.client_status.cpu_is_x64 ? "x64 " : "");
        tooltip += (row.client_status.cpu_is_vm ? "VM" : "");
        tooltip += '\n';
        tooltip += "CPU Cache L2/L3: " + cache(row.client_status.cpu_l2) + " MB/" + cache(row.client_status.cpu_l3) + " MB";
        tooltip += '\n';
        tooltip += "CPU Nodes: " + row.client_status.cpu_nodes;
        tooltip += '\n';
        tooltip += "Max CPU usage: " + (row.client_status.max_cpu_usage > 0 ? row.client_status.max_cpu_usage : "100") + "%";
        tooltip += '\n';
        tooltip += "Huge Pages: " + (row.client_status.hugepages_available ? " available, " : " unavailable, ");
        tooltip += (row.client_status.hugepages_enabled ? "enabled (" + row.client_status.total_hugepages + "/" + row.client_status.total_pages + ")" : "disabled");
        tooltip += '\n';
        tooltip += "Used Threads: " + row.client_status.current_threads;
        tooltip += (row.client_status.hash_factor > 1 ? " [" + row.client_status.hash_factor + "x multi hash mode]" : "");
        tooltip += '\n';
        tooltip += "Memory Free/Total: " + memory(row.client_status.free_memory) + " GB/" + memory(row.client_status.total_memory) + " GB";
        tooltip += '\n';

        if (row.client_status.gpu_info_list) {
            for (var id in row.client_status.gpu_info_list) {
                tooltip += "GPU #" + row.client_status.gpu_info_list[id].gpu_info.device_idx + ": ";
                tooltip += row.client_status.gpu_info_list[id].gpu_info.name + ", "
                tooltip += "intensity: " + row.client_status.gpu_info_list[id].gpu_info.raw_intensity + " ";
                tooltip += "(" + row.client_status.gpu_info_list[id].gpu_info.work_size + "/" + row.client_status.gpu_info_list[id].gpu_info.max_work_size + "), ";
                tooltip += "cu: " + row.client_status.gpu_info_list[id].gpu_info.compute_units;
                tooltip += '\n';
            }
        }

        tooltip += "Client IP: " + row.client_status.external_ip;
        tooltip += '\n';
        tooltip += "Version: " + row.client_status.version;
        tooltip += '\n';
        tooltip += "Status: " + online ? "Online" : "Offline";

        if (online) {
            return `<span data-toggle="tooltip" title="${tooltip}"><div class="online">${data}</div></span>`;
        } else {
            return `<span data-toggle="tooltip" title="${tooltip}"><div class="offline">${data}</div></span>`;
        }
    }

    return data;
}



function memory(data) {
    return Math.round(data / 1024 / 1024 / 1024 * 10) / 10;
}

function cache(data) {
    return Math.round(data / 1024 * 100) / 100;
}

function isOnline(lastStatus, TRESHOLD_IN_MS, currentServerTime) {
    var threshold = currentServerTime - TRESHOLD_IN_MS;
    if (lastStatus > threshold) {
        return true;
    } else {
        return false;
    }
}

function version(data, latestVersion, row) {
    var clientVersion = parseInt(row.client_status.version.split('.').join(""));

    if (clientVersion < 1000) {
        clientVersion = clientVersion * 10;
    }

    if (latestVersion > clientVersion) {
        return `<span data-toggle="tooltip" title="Outdated"><div class="offline">${data}</div></span>`;
    } else {
        return data;
    }
}
function algoAndPowVariantName(data, type, row) {
    var algo = row.client_status.current_algo_name;
    var powVariant = row.client_status.current_pow_variant_name;

    if (powVariant !== "") {
        return algo + " / " + powVariant
    } else {
        return algo;
    }
}

function uptime(data, type, row) {
    if (type !== 'sort') {
        var lastStatus = row.client_status.last_status_update * 1000;

        if (isOnline(lastStatus)) {
            return numeral(data / 1000).format('00:00:00:00');
        } else {
            return "";
        }
    }

    return data;
}

function laststatus() {
    return function (data, type) {
        if (type !== 'sort') {
            var date = new Date(data * 1000);
            return `<span data-toggle="tooltip" title="${date}">${$.timeago(date)}</span>`;
        }

        return data;
    };
}
function round(data) {
    return Math.round(data * 100) / 100;
}

export { version, algoAndPowVariantName, memory, cache, isOnline, uptime, laststatus, clientInfo, clientStatus, round };