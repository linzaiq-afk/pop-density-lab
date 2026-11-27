(() => {
  "use strict";

  // При открытии index.html напрямую (file://) API ищется на локальном порту Go.
  const API_BASE = window.location.protocol === "file:" ? "http://localhost:8081" : "";
  const DEFAULT_SELECTION = ["msk", "kzn", "tyo"];
  const PALETTE = [
    "#2563eb", "#16a34a", "#dc2626", "#d97706", "#7c3aed",
    "#0891b2", "#db2777", "#65a30d", "#ea580c", "#4f46e5",
  ];

  const state = {
    entities: [],
    selected: new Set(DEFAULT_SELECTION),
    type: "all",
    query: "",
    logScale: false,
  };

  const fmtInt = new Intl.NumberFormat("ru-RU", { maximumFractionDigits: 0 });
  const fmtNum = new Intl.NumberFormat("ru-RU", { maximumFractionDigits: 2 });

  const el = {
    apiStatus: document.getElementById("apiStatus"),
    apiStatusWrap: document.getElementById("apiStatusWrap"),
    errorBanner: document.getElementById("errorBanner"),
    kpiEntities: document.getElementById("kpiEntities"),
    kpiEntitiesNote: document.getElementById("kpiEntitiesNote"),
    kpiLeader: document.getElementById("kpiLeader"),
    kpiLeaderNote: document.getElementById("kpiLeaderNote"),
    kpiGrowth: document.getElementById("kpiGrowth"),
    kpiGrowthNote: document.getElementById("kpiGrowthNote"),
    kpiSelected: document.getElementById("kpiSelected"),
    kpiSelectedNote: document.getElementById("kpiSelectedNote"),
    typeFilter: document.getElementById("typeFilter"),
    searchInput: document.getElementById("searchInput"),
    logScale: document.getElementById("logScale"),
    checkboxes: document.getElementById("entityCheckboxes"),
    resetBtn: document.getElementById("resetBtn"),
    selectAllBtn: document.getElementById("selectAllBtn"),
    csvBtn: document.getElementById("csvBtn"),
    canvas: document.getElementById("densityChart"),
    chartEmpty: document.getElementById("chartEmpty"),
    selectionSummary: document.getElementById("selectionSummary"),
    tbody: document.querySelector("#statsTable tbody"),
    tableCount: document.getElementById("tableCount"),
    calcPreset: document.getElementById("calcPreset"),
    calcPopulation: document.getElementById("calcPopulation"),
    calcArea: document.getElementById("calcArea"),
    calcResult: document.getElementById("calcResult"),
    calcNote: document.getElementById("calcNote"),
  };

  let chart = null;

  // ------------------------------------------------------------------ utils

  function typeLabel(type) {
    return type === "city" ? "Город" : "Страна";
  }

  function statFor(entity, year) {
    return entity.history.find((stat) => stat.year === year) || null;
  }

  function formatDensity(value) {
    return value === null || value === undefined ? "—" : fmtNum.format(value);
  }

  function colorFor(entity) {
    const index = state.entities.findIndex((item) => item.id === entity.id);
    return PALETTE[(index >= 0 ? index : 0) % PALETTE.length];
  }

  async function api(path) {
    const response = await fetch(API_BASE + path, { headers: { Accept: "application/json" } });
    if (!response.ok) {
      throw new Error(`${response.status} ${response.statusText}`);
    }
    return response.json();
  }

  function filteredEntities() {
    const query = state.query.trim().toLowerCase();
    return state.entities.filter((entity) => {
      if (state.type !== "all" && entity.type !== state.type) return false;
      if (!query) return true;
      return (
        entity.name.toLowerCase().includes(query) ||
        entity.country.toLowerCase().includes(query)
      );
    });
  }

  function growthPercent(entity) {
    if (!entity.history.length) return null;
    const first = entity.history[0];
    const last = entity.history[entity.history.length - 1];
    if (!first.density) return null;
    return ((last.density - first.density) / first.density) * 100;
  }

  // ------------------------------------------------------------------ chips

  function renderChips() {
    const entities = filteredEntities();
    el.checkboxes.replaceChildren();

    if (!entities.length) {
      const empty = document.createElement("p");
      empty.className = "hint";
      empty.textContent = "Ничего не найдено — измените фильтр или поисковый запрос.";
      el.checkboxes.append(empty);
      return;
    }

    for (const entity of entities) {
      const label = document.createElement("label");
      label.className = "chip";

      const checkbox = document.createElement("input");
      checkbox.type = "checkbox";
      checkbox.checked = state.selected.has(entity.id);
      checkbox.addEventListener("change", () => {
        if (checkbox.checked) {
          state.selected.add(entity.id);
        } else {
          state.selected.delete(entity.id);
        }
        renderChart();
        renderTable();
        renderKpis();
      });

      const name = document.createElement("span");
      name.textContent = entity.name;

      const meta = document.createElement("small");
      meta.textContent = entity.type === "city" ? entity.country : "страна";

      label.append(checkbox, name, meta);
      el.checkboxes.append(label);
    }
  }

  // ------------------------------------------------------------------ chart

  function renderChart() {
    const entities = state.entities.filter((entity) => state.selected.has(entity.id));

    el.selectionSummary.textContent = entities.length
      ? `На графике: ${entities.map((entity) => entity.name).join(", ")}`
      : "На графике пока ничего не выбрано";

    if (!entities.length) {
      if (chart) {
        chart.destroy();
        chart = null;
      }
      el.chartEmpty.hidden = false;
      return;
    }

    el.chartEmpty.hidden = true;

    const years = [...new Set(entities.flatMap((entity) => entity.history.map((stat) => stat.year)))].sort(
      (a, b) => a - b,
    );

    const datasets = entities.map((entity) => {
      const color = colorFor(entity);
      const byYear = new Map(entity.history.map((stat) => [stat.year, stat]));
      return {
        label: entity.name,
        entityId: entity.id,
        data: years.map((year) => (byYear.has(year) ? byYear.get(year).density : null)),
        borderColor: color,
        backgroundColor: `${color}1a`,
        borderWidth: 2.5,
        pointRadius: 3,
        pointHoverRadius: 5,
        pointBackgroundColor: "#ffffff",
        pointBorderColor: color,
        pointBorderWidth: 2,
        tension: 0.25,
        spanGaps: true,
        fill: false,
      };
    });

    if (chart) {
      chart.destroy();
    }

    chart = new Chart(el.canvas.getContext("2d"), {
      type: "line",
      data: { labels: years, datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: { mode: "index", intersect: false },
        plugins: {
          legend: {
            position: "bottom",
            labels: { usePointStyle: true, boxWidth: 8, padding: 16, font: { size: 12 } },
          },
          tooltip: {
            padding: 12,
            callbacks: {
              title: (items) => `Год: ${items[0].label}`,
              label: (context) => {
                if (context.parsed.y === null || context.parsed.y === undefined) {
                  return `${context.dataset.label}: нет данных`;
                }
                return `${context.dataset.label}: ${fmtNum.format(context.parsed.y)} чел./км²`;
              },
              afterLabel: (context) => {
                const entity = state.entities.find((item) => item.id === context.dataset.entityId);
                if (!entity) return "";
                const stat = entity.history.find((item) => item.year === years[context.dataIndex]);
                if (!stat) return "";
                return [
                  `Население: ${fmtInt.format(stat.population)} чел.`,
                  `Площадь: ${fmtNum.format(stat.area_km2)} км²`,
                ];
              },
            },
          },
        },
        scales: {
          y: {
            type: state.logScale ? "logarithmic" : "linear",
            beginAtZero: !state.logScale,
            title: { display: true, text: "Плотность, чел./км²" },
            grid: { color: "rgba(148, 163, 184, 0.18)" },
          },
          x: {
            title: { display: true, text: "Год" },
            grid: { display: false },
          },
        },
      },
    });
  }

  // ------------------------------------------------------------------ table

  function renderTable() {
    const entities = filteredEntities()
      .slice()
      .sort((a, b) => {
        const densityA = statFor(a, 2026)?.density ?? 0;
        const densityB = statFor(b, 2026)?.density ?? 0;
        return densityB - densityA;
      });

    el.tbody.replaceChildren();
    el.tableCount.textContent = `${entities.length} из ${state.entities.length}`;

    for (const entity of entities) {
      const growth = growthPercent(entity);
      const row = document.createElement("tr");
      row.dataset.id = entity.id;
      row.classList.toggle("is-selected", state.selected.has(entity.id));
      row.tabIndex = 0;
      row.setAttribute("role", "button");
      row.setAttribute("aria-pressed", String(state.selected.has(entity.id)));
      row.title = "Добавить или убрать объект с графика";

      row.append(
        cell(entity.name, "strong"),
        cell(typeLabel(entity.type)),
        cell(entity.type === "city" ? entity.country : "—"),
        cell(formatDensity(statFor(entity, 1960)?.density), "num"),
        cell(formatDensity(statFor(entity, 2000)?.density), "num"),
        cell(formatDensity(statFor(entity, 2026)?.density), "num"),
        cell(
          growth === null ? "—" : `${growth >= 0 ? "+" : ""}${fmtNum.format(growth)} %`,
          growth === null ? "num" : growth >= 0 ? "num pos" : "num neg",
        ),
      );

      const toggle = () => {
        if (state.selected.has(entity.id)) {
          state.selected.delete(entity.id);
        } else {
          state.selected.add(entity.id);
        }
        renderChips();
        renderChart();
        renderTable();
        renderKpis();
      };

      row.addEventListener("click", toggle);
      row.addEventListener("keydown", (event) => {
        if (event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          toggle();
        }
      });

      el.tbody.append(row);
    }
  }

  function cell(text, className) {
    const td = document.createElement("td");
    td.textContent = text;
    if (className) td.className = className;
    return td;
  }

  // -------------------------------------------------------------------- kpi

  function renderKpis() {
    const entities = filteredEntities();
    el.kpiEntities.textContent = fmtInt.format(entities.length);
    el.kpiEntitiesNote.textContent = `из ${fmtInt.format(state.entities.length)} в базе`;

    let leader = null;
    for (const entity of entities) {
      const stat = statFor(entity, 2026);
      if (stat && (!leader || stat.density > leader.density)) {
        leader = { name: entity.name, density: stat.density };
      }
    }
    el.kpiLeader.textContent = leader ? leader.name : "—";
    el.kpiLeaderNote.textContent = leader
      ? `${fmtInt.format(leader.density)} чел./км² в 2026`
      : "нет данных";

    const growths = entities
      .map((entity) => growthPercent(entity))
      .filter((value) => value !== null);
    const average = growths.length
      ? growths.reduce((sum, value) => sum + value, 0) / growths.length
      : null;
    el.kpiGrowth.textContent =
      average === null ? "—" : `${average >= 0 ? "+" : ""}${fmtNum.format(average)} %`;
    el.kpiGrowthNote.textContent = "средний рост плотности";

    el.kpiSelected.textContent = fmtInt.format(state.selected.size);
    el.kpiSelectedNote.textContent = "объектов на графике";
  }

  // ------------------------------------------------------------- calculator

  function updateCalculator() {
    const population = Number.parseFloat(el.calcPopulation.value);
    const area = Number.parseFloat(el.calcArea.value);

    if (!Number.isFinite(population) || !Number.isFinite(area) || area <= 0 || population < 0) {
      el.calcResult.textContent = "—";
      el.calcNote.textContent = "Укажите неотрицательное население и площадь > 0.";
      return;
    }

    el.calcResult.textContent = fmtNum.format(population / area);
    el.calcNote.textContent = `${fmtInt.format(population)} чел. ÷ ${fmtNum.format(area)} км²`;
  }

  function applyCalculatorPreset(id) {
    const entity = state.entities.find((item) => item.id === id);
    const stat = entity && statFor(entity, 2026);
    if (!stat) return;
    el.calcPopulation.value = String(stat.population);
    el.calcArea.value = String(stat.area_km2);
    updateCalculator();
  }

  function initCalculator() {
    const entities = state.entities.slice().sort((a, b) => a.name.localeCompare(b.name, "ru"));
    el.calcPreset.replaceChildren();

    for (const entity of entities) {
      const option = document.createElement("option");
      option.value = entity.id;
      option.textContent = `${entity.name} — ${typeLabel(entity.type).toLowerCase()}`;
      el.calcPreset.append(option);
    }

    if (entities.some((entity) => entity.id === "msk")) {
      el.calcPreset.value = "msk";
    }
    applyCalculatorPreset(el.calcPreset.value);
  }

  // ------------------------------------------------------------------- csv

  function exportCsv() {
    const entities = filteredEntities()
      .slice()
      .sort((a, b) => (statFor(b, 2026)?.density ?? 0) - (statFor(a, 2026)?.density ?? 0));

    const rows = [
      ["Объект", "Категория", "Страна", "1960, чел./км²", "2000, чел./км²", "2026, чел./км²", "Изменение 1960→2026, %"],
    ];

    for (const entity of entities) {
      const growth = growthPercent(entity);
      rows.push([
        entity.name,
        typeLabel(entity.type),
        entity.country,
        formatDensity(statFor(entity, 1960)?.density),
        formatDensity(statFor(entity, 2000)?.density),
        formatDensity(statFor(entity, 2026)?.density),
        growth === null ? "—" : fmtNum.format(growth),
      ]);
    }

    const csv = rows
      .map((row) => row.map((value) => `"${String(value).replaceAll('"', '""')}"`).join(";"))
      .join("\r\n");

    // BOM нужен, чтобы Excel корректно открыл кириллицу.
    const blob = new Blob(["\uFEFF", csv], { type: "text/csv;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "population-density.csv";
    document.body.append(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  }

  // ---------------------------------------------------------------- events

  function renderAll() {
    renderChips();
    renderChart();
    renderTable();
    renderKpis();
  }

  function resetAll() {
    const available = new Set(state.entities.map((entity) => entity.id));
    state.selected = new Set(DEFAULT_SELECTION.filter((id) => available.has(id)));
    if (!state.selected.size && state.entities[0]) {
      state.selected.add(state.entities[0].id);
    }
    state.type = "all";
    state.query = "";
    state.logScale = false;

    el.typeFilter.value = "all";
    el.searchInput.value = "";
    el.logScale.checked = false;

    renderAll();
  }

  function bindEvents() {
    el.typeFilter.addEventListener("change", () => {
      state.type = el.typeFilter.value;
      renderAll();
    });

    el.searchInput.addEventListener("input", () => {
      state.query = el.searchInput.value;
      renderAll();
    });

    el.logScale.addEventListener("change", () => {
      state.logScale = el.logScale.checked;
      renderChart();
    });

    el.selectAllBtn.addEventListener("click", () => {
      for (const entity of filteredEntities()) {
        state.selected.add(entity.id);
      }
      renderAll();
    });

    el.resetBtn.addEventListener("click", resetAll);
    el.csvBtn.addEventListener("click", exportCsv);

    el.calcPreset.addEventListener("change", () => applyCalculatorPreset(el.calcPreset.value));
    el.calcPopulation.addEventListener("input", updateCalculator);
    el.calcArea.addEventListener("input", updateCalculator);
  }

  // ------------------------------------------------------------------- init

  async function init() {
    bindEvents();
    updateCalculator();

    try {
      const entities = await api("/api/entities");
      state.entities = entities;

      const available = new Set(entities.map((entity) => entity.id));
      state.selected = new Set(DEFAULT_SELECTION.filter((id) => available.has(id)));
      if (!state.selected.size && entities[0]) {
        state.selected.add(entities[0].id);
      }

      el.apiStatus.textContent = `API: онлайн · ${entities.length} объектов`;
      el.apiStatusWrap.classList.add("is-online");

      initCalculator();
      renderAll();
    } catch (error) {
      console.error(error);
      el.apiStatus.textContent = "API: ошибка";
      el.apiStatusWrap.classList.add("is-error");
      el.errorBanner.hidden = false;
      el.errorBanner.textContent =
        `Не удалось загрузить данные (${error.message}). ` +
        "Запустите стенд командой docker compose up -d --build и обновите страницу.";
      el.chartEmpty.hidden = false;
    }
  }

  init();
})();
