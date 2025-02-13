<template>
  <div style="height: 30vh; margin: 0 auto">
    <canvas id="hashrateChart"></canvas>
    <br /><br />
    <canvas id="minerChart"></canvas>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue';
import Chart from 'chart.js';
import 'chartjs-plugin-colorschemes';
import axios from 'axios';

export default {
  props: {
    CHART_RELOAD_INTERVAL_IN_MS: {
      type: Number,
      required: true,
    },
  },
  setup(props) {
    const hashrateChart = ref(null);
    const minerChart = ref(null);

    onMounted(() => {
      hashrateChart.value = createChart('hashrateChart', 'Hashrates', 'Hashrate');
      minerChart.value = createChart('minerChart', 'Miners', 'Miner', true);
      updateCharts();

      setInterval(() => {
        updateCharts();
      }, props.CHART_RELOAD_INTERVAL_IN_MS);
    });

    const createChart = (canvasId, title, labelString, integerValues = false) => {
      return new Chart(document.getElementById(canvasId).getContext('2d'), {
        type: 'line',
        data: { datasets: [] },
        options: {
          title: { display: true, text: title },
          tooltips: { mode: 'nearest', intersect: false },
          scales: {
            yAxes: [
              {
                id: canvasId.replace('Chart', ''), // Dynamische ID für y-Achse
                ticks: {
                  beginAtZero: true,
                  callback: integerValues
                    ? function (value) {
                        if (value % 1 === 0) {
                          return value;
                        }
                      }
                    : null,
                },
                scaleLabel: { display: true, labelString: labelString, fontSize: 16 },
              },
            ],
            xAxes: [
              {
                type: 'time',
                distribution: 'series',
                time: { unit: 'minute', displayFormats: { minute: 'HH:mm' } },
                ticks: { stepSize: 10, autoSkip: true, maxTicksLimit: 60 },
              },
            ],
          },
          elements: { point: { radius: 0 } },
          responsive: true,
          maintainAspectRatio: false,
          plugins: { colorschemes: { scheme: 'brewer.Paired12' } },
        },
      });
    };

    const updateChartData = (chart, results, dataLabel) => {
      var datasets = results['client_statistics'].map(function (entry) {
        return {
          label: entry.algo,
          type: 'line',
          data: entry.statistics.map(function (stat) {
            return { x: new Date(stat.timestamp), y: stat[dataLabel] };
          }),
          spanGraphs: false,
          fill: dataLabel === 'hashrate', // Nur für Hashrate-Chart füllen
        };
      });

      chart.data.datasets = datasets;
      chart.update();
    };

    const updateCharts = async () => {
      try {
        const response = await axios.get('/admin/getClientStatistics');
        const results = response.data;

        var algos = results["client_statistics"].map(function(e) {
            return e.algo;
        });

        var statistics = results["client_statistics"].map(function(e) {
            return e.statistics;
        });

        var hashrateDatasets=[], minerDatasets=[];
        for(var j = 0; j < results["client_statistics"].length; j++) {
            var hashrates = statistics[j].map(function(e) {
                return {
                    x: new Date(e.timestamp),
                    y: e.hashrate
                }
            });

            var miners = statistics[j].map(function(e) {
                return {
                    x: new Date(e.timestamp),
                    y: e.miner
                }
            });

            hashrateDatasets.push({label: algos[j], type: 'line', data: hashrates, spanGraphs: false, fill: true});
            minerDatasets.push({label: algos[j], type: 'line', data: miners, spanGraphs: false, fill: false});
        }

        hashrateChart.value.data.datasets = hashrateDatasets;
        hashrateChart.value.update();

        minerChart.value.data.datasets = minerDatasets;
        minerChart.value.update();

      } catch (error) {
        console.error('Error fetching chart data:', error);
      }
    };

    return {
      hashrateChart,
      minerChart,
      createChart,
      updateChartData,
      updateCharts,
    };
  },
};
</script>