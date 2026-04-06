import { readFile } from "fs/promises";

const BASE_URL = "http://localhost:3000";

const credentials = {
  email: "admin@gmail.com",
  password: "admin",
};

async function loadRecords() {
  const data = await readFile("./records.json", "utf-8");
  return JSON.parse(data);
}

async function login() {
  const res = await fetch(`${BASE_URL}/users/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(`Login failed: ${err}`);
  }

  const data = await res.json();

  const token = data.access_token;

  if (!token) {
    throw new Error("No access_token in response");
  }

  return token;
}

async function insertRecord(record, token) {
  const res = await fetch(`${BASE_URL}/records`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(record),
  });

  if (!res.ok) {
    const err = await res.text();
    console.error("❌ Failed:", record, err);
    return;
  }

  console.log(
    `✅ Inserted: ${record.type} | ${record.category} | ${record.amount}`,
  );
}

async function seed() {
  try {
    console.log("📂 Loading records...");
    const records = await loadRecords();
    console.log(`✅ Loaded ${records.length} records`);

    console.log("🔐 Logging in...");
    const token = await login();
    console.log("✅ Got token");

    for (let i = 0; i < records.length; i++) {
      await insertRecord(records[i], token);
    }

    console.log("🎉 Seeding complete");
  } catch (err) {
    console.error("🔥 Error:", err.message);
  }
}

seed();
