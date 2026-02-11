const API_BASE = "/api/v1/comments";

const state = {
  page: 1,
  limit: 10,
  sort: "created_at_desc",
  searching: false,
  searchQuery: "",
};

const commentsRoot = document.getElementById("commentsRoot");
const searchInput = document.getElementById("searchInput");
const searchBtn = document.getElementById("searchBtn");
const sortSelect = document.getElementById("sortSelect");
const prevBtn = document.getElementById("prevBtn");
const nextBtn = document.getElementById("nextBtn");
const pageInfo = document.getElementById("pageInfo");
const limitSelect = document.getElementById("limitSelect");
const authorInput = document.getElementById("authorInput");
const contentInput = document.getElementById("contentInput");
const createTopBtn = document.getElementById("createTopBtn");
const msg = document.getElementById("msg");

function showMessage(text, isError) {
  msg.textContent = text || "";
  msg.className = isError ? "error small" : "small";
  if (text)
    setTimeout(() => {
      msg.textContent = "";
    }, 4000);
}

async function fetchJSON(url, opts = {}) {
  const res = await fetch(url, opts);
  const text = await res.text();
  let data = null;
  try {
    data = text ? JSON.parse(text) : null;
  } catch (e) {
    throw new Error("Invalid JSON response");
  }

  if (!res.ok) {
    const errMsg =
      data && data.error ? data.error : res.status + " " + res.statusText;
    throw new Error(errMsg);
  }

  return data && data.result !== undefined ? data.result : data;
}

async function loadComments() {
  commentsRoot.innerHTML = '<div class="small">Loading...</div>';

  try {
    if (state.searching && state.searchQuery.trim() !== "") {
      const limit = 100;
      const url = `${API_BASE}?page=1&limit=${limit}&sort=${state.sort}`;
      const roots = await fetchJSON(url);
      const filtered = filterForestByQuery(roots, state.searchQuery);
      renderComments(filtered);
      pageInfo.textContent = `Search: "${state.searchQuery}" — ${countNodes(filtered)} results`;
    } else {
      const url = `${API_BASE}?page=${state.page}&limit=${state.limit}&sort=${state.sort}`;
      const roots = await fetchJSON(url);
      renderComments(roots);
      pageInfo.textContent = `Page ${state.page}`;
    }
  } catch (err) {
    commentsRoot.innerHTML = "";
    showMessage(err.message, true);
  }
}

function countNodes(list) {
  let cnt = 0;
  function walk(node) {
    cnt++;
    (node.children || []).forEach(walk);
  }
  list.forEach(walk);
  return cnt;
}

function filterForestByQuery(list, q) {
  const ql = q.trim().toLowerCase();
  if (!ql) return list;

  function filterNode(node) {
    const me =
      (node.content || "").toLowerCase().includes(ql) ||
      (node.author || "").toLowerCase().includes(ql);

    const children = (node.children || []).map(filterNode).filter(Boolean);

    if (me || children.length) {
      return { ...node, children };
    }
    return null;
  }

  return list.map(filterNode).filter(Boolean);
}

function renderComments(list) {
  commentsRoot.innerHTML = "";

  if (!list || list.length === 0) {
    commentsRoot.innerHTML = '<div class="small">No comments</div>';
    return;
  }

  list.forEach((c) => {
    const el = renderNode(c, 0);
    commentsRoot.appendChild(el);
  });
}

function renderNode(node, level) {
  const wrap = document.createElement("div");
  wrap.style.marginLeft = level * 16 + "px";
  wrap.className = "comment";

  const meta = document.createElement("div");
  meta.className = "meta";
  meta.innerHTML =
    `<strong>${escapeHtml(node.author || "—")}</strong> · ` +
    `<span>${formatDate(node.created_at)}</span>`;
  wrap.appendChild(meta);

  const content = document.createElement("div");
  content.className = "content";
  content.textContent = node.content || "";
  wrap.appendChild(content);

  const actions = document.createElement("div");
  actions.className = "actions";

  const replyBtn = document.createElement("button");
  replyBtn.textContent = "Reply";
  replyBtn.onclick = () => openReplyForm(wrap, node.id);
  actions.appendChild(replyBtn);

  const viewBtn = document.createElement("button");
  viewBtn.className = "ghost";
  viewBtn.textContent = "Open thread";
  viewBtn.onclick = () => openThread(node.id);
  actions.appendChild(viewBtn);

  const delBtn = document.createElement("button");
  delBtn.className = "ghost";
  delBtn.textContent = "Delete";
  delBtn.onclick = async () => {
    if (!confirm("Delete this comment and all nested replies?")) return;
    try {
      await fetchJSON(`${API_BASE}/${node.id}`, { method: "DELETE" });
      showMessage("Deleted");
      await loadComments();
    } catch (err) {
      showMessage(err.message, true);
    }
  };
  actions.appendChild(delBtn);

  wrap.appendChild(actions);

  const children = node.children || [];
  if (children.length) {
    const list = document.createElement("div");
    list.style.marginTop = "8px";
    children.forEach((ch) => {
      list.appendChild(renderNode(ch, level + 1));
    });
    wrap.appendChild(list);
  }

  return wrap;
}

function openReplyForm(container, parentID) {
  if (container.querySelector(".reply-form")) return;

  const form = document.createElement("div");
  form.className = "reply-form";

  const author = document.createElement("input");
  author.placeholder = "Name";

  const ta = document.createElement("textarea");
  ta.rows = 2;
  ta.placeholder = "Reply text";
  ta.style.flex = "1";

  const send = document.createElement("button");
  send.textContent = "Send";

  const cancel = document.createElement("button");
  cancel.className = "ghost";
  cancel.textContent = "Cancel";

  send.onclick = async () => {
    const payload = {
      parent_id: parentID,
      author: author.value,
      content: ta.value,
    };
    try {
      await fetchJSON(API_BASE, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      showMessage("Reply added");
      await loadComments();
    } catch (err) {
      showMessage(err.message, true);
    }
  };

  cancel.onclick = () => form.remove();

  form.appendChild(author);
  form.appendChild(ta);
  form.appendChild(send);
  form.appendChild(cancel);
  container.appendChild(form);
}

async function openThread(id) {
  commentsRoot.innerHTML = '<div class="small">Loading thread...</div>';
  try {
    const url = `${API_BASE}?parent=${id}`;
    const roots = await fetchJSON(url);
    renderComments(roots);
    pageInfo.textContent = `Thread ${id}`;
  } catch (err) {
    showMessage(err.message, true);
  }
}

async function createTopComment() {
  const author = authorInput.value.trim();
  const content = contentInput.value.trim();

  if (!author) {
    showMessage("Author is required", true);
    return;
  }

  if (!content) {
    showMessage("Content is required", true);
    return;
  }

  const payload = { parent_id: null, author, content };

  try {
    await fetchJSON(API_BASE, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    authorInput.value = "";
    contentInput.value = "";
    showMessage("Comment added");
    await loadComments();
  } catch (err) {
    showMessage(err.message, true);
  }
}

function formatDate(iso) {
  try {
    if (!iso) return "";
    const d = new Date(iso);
    return d.toLocaleString();
  } catch (e) {
    return iso;
  }
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, function (ch) {
    return {
      "&": "&amp;",
      "<": "&lt;",
      ">": "&gt;",
      '"': "&quot;",
      "'": "&#39;",
    }[ch];
  });
}

searchBtn.addEventListener("click", () => {
  const q = searchInput.value.trim();
  state.searching = !!q;
  state.searchQuery = q;
  loadComments();
});

sortSelect.addEventListener("change", () => {
  state.sort = sortSelect.value;
  loadComments();
});

limitSelect.addEventListener("change", () => {
  state.limit = parseInt(limitSelect.value, 10);
  state.page = 1;
  loadComments();
});

prevBtn.addEventListener("click", () => {
  if (state.page > 1) {
    state.page--;
    loadComments();
  }
});

nextBtn.addEventListener("click", () => {
  state.page++;
  loadComments();
});

createTopBtn.addEventListener("click", createTopComment);

loadComments();
