"use strict";

const $ = (id) => document.getElementById(id);

function spec() {
  return {
    metal: $("metal").value,
    i_corr: Number($("icorr").value),
    M: Number($("molar").value),
    n: Number($("valence").value),
    rho: Number($("rho").value),
    area: Number($("area").value),
    duration_y: Number($("years").value),
  };
}

function showError(msg) {
  $("err-text").textContent = msg;
  $("err-panel").hidden = false;
}

function hideError() {
  $("err-panel").hidden = true;
}

function fillTable(rows) {
  const tb = $("out-table");
  tb.innerHTML = "";
  for (const [k, v] of rows) {
    const tr = document.createElement("tr");
    tr.innerHTML = "<td>" + k + "</td><td>" + v + "</td>";
    tb.appendChild(tr);
  }
}

async function loadExample() {
  hideError();
  $("hint").textContent = "正在读取 /api/example …";
  try {
    const resp = await fetch("/api/example");
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "示例加载失败");
    $("metal").value = data.metal;
    $("icorr").value = data.i_corr;
    $("molar").value = data.M;
    $("valence").value = data.n;
    $("rho").value = data.rho;
    $("area").value = data.area;
    $("years").value = data.duration_y;
    $("hint").textContent = "已加载铁海水算例。点“求速率”看结果。";
  } catch (e) {
    showError(String(e));
  }
}

async function runRate() {
  hideError();
  try {
    const resp = await fetch("/api/rate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(spec()),
    });
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "HTTP " + resp.status);
    $("out-panel").hidden = false;
    fillTable([
      ["i_corr (μA/cm²)", data.i_corr.toPrecision(5)],
      ["腐蚀深度 CR (mm/y)", data.corrosion_rate_mm_y.toPrecision(5)],
      ["腐蚀深度 (μm/y)", data.corrosion_rate_um_y.toPrecision(5)],
      ["年失重 (g/cm²·y)", data.annual_mass_loss.toPrecision(5)],
      ["累计失重 (g)", data.cumulative_mass_g.toPrecision(5)],
    ]);
    $("hint").textContent = "数值全部来自后端法拉第公式。";
  } catch (e) {
    showError(String(e));
  }
}

$("btn-example").addEventListener("click", loadExample);
$("btn-rate").addEventListener("click", runRate);
