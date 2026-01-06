#!/usr/bin/env python3
import datetime
import json
import os
import sys
import time
import urllib.error
import urllib.request


BASE_URL = os.getenv("BASE_URL", "http://localhost:8080").rstrip("/")
ADMIN_USERNAME = os.getenv("ADMIN_USERNAME", "admin")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "admin123")
MEMBER_PASSWORD = os.getenv("MEMBER_PASSWORD", "member123")
KEEP_DATA = os.getenv("KEEP_DATA", "").lower() in {"1", "true", "yes"}


def fail(message):
    print(f"FAIL: {message}")
    sys.exit(1)


def parse_body(raw):
    if not raw:
        return None
    text = raw.decode("utf-8", errors="replace")
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        return text


def get_field(obj, *keys):
    if not isinstance(obj, dict):
        return None
    for key in keys:
        if key in obj:
            return obj[key]

    def normalize_key(value):
        return "".join(ch for ch in value.lower() if ch.isalnum())

    normalized_map = {normalize_key(k): v for k, v in obj.items()}
    for key in keys:
        normalized = normalize_key(key)
        if normalized in normalized_map:
            return normalized_map[normalized]
    return None


def request(method, path, token=None, json_body=None, expected_status=None):
    url = f"{BASE_URL}{path}"
    data = None
    headers = {}

    if json_body is not None:
        data = json.dumps(json_body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    if token:
        headers["Authorization"] = f"Bearer {token}"

    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req) as resp:
            status = resp.getcode()
            body = parse_body(resp.read())
    except urllib.error.HTTPError as err:
        status = err.code
        body = parse_body(err.read())
    except urllib.error.URLError as err:
        fail(f"{method} {path} -> {err}")

    if expected_status is not None and status != expected_status:
        fail(f"{method} {path} expected {expected_status}, got {status}: {body}")

    return status, body


def to_rfc3339_z(dt):
    return dt.astimezone(datetime.UTC).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def login(username, password):
    _, body = request(
        "POST",
        "/api/auth/login",
        json_body={"username": username, "password": password},
        expected_status=200,
    )
    if not isinstance(body, dict) or "token" not in body:
        fail("login response missing token")
    return body["token"], body.get("role")


def main():
    print(f"Testing API at {BASE_URL}")

    admin_token, admin_role = login(ADMIN_USERNAME, ADMIN_PASSWORD)
    if admin_role != "admin":
        fail(f"admin login role mismatch: {admin_role}")
    print("OK: admin login")

    _, accounts = request(
        "GET",
        "/api/admin/accounts",
        token=admin_token,
        expected_status=200,
    )
    if not isinstance(accounts, list):
        fail("accounts list response is not a list")
    admin_account = next(
        (acc for acc in accounts if get_field(acc, "username") == ADMIN_USERNAME), None
    )
    if not admin_account:
        fail("admin account not found in list")
    admin_id = get_field(admin_account, "id")
    if admin_id is None:
        fail("admin account missing id")
    print("OK: list admin accounts")

    suffix = int(time.time())
    member_username = f"member_{suffix}"
    member_email = f"{member_username}@example.com"
    member_phone = int(f"88{suffix % 1000000:06d}")

    _, member_account = request(
        "POST",
        "/api/admin/accounts",
        token=admin_token,
        json_body={
            "username": member_username,
            "email": member_email,
            "phone": member_phone,
            "password": MEMBER_PASSWORD,
        },
        expected_status=201,
    )
    if not isinstance(member_account, dict) or "id" not in member_account:
        member_id = get_field(member_account, "id")
        if member_id is None:
            fail("member account create response missing id")
    else:
        member_id = member_account["id"]
    print("OK: create member account")

    now = datetime.datetime.now(datetime.UTC)
    start_date = to_rfc3339_z(now)
    end_date = to_rfc3339_z(now + datetime.timedelta(days=1))

    _, log_head = request(
        "POST",
        "/api/admin/log-heads",
        token=admin_token,
        json_body={
            "subject": f"Test Subject {suffix}",
            "start_date": start_date,
            "end_date": end_date,
            "writer_id_list": [member_id],
            "owner_id": member_id,
        },
        expected_status=201,
    )
    log_head_id = get_field(log_head, "id")
    if not isinstance(log_head, dict) or log_head_id is None:
        fail("log head create response missing id")
    print("OK: create log head")

    request("GET", "/api/log-heads", token=admin_token, expected_status=200)
    print("OK: list log heads")

    request("GET", "/api/log-heads/writable", token=admin_token, expected_status=200)
    print("OK: list writable log heads (admin)")

    member_token, member_role = login(member_username, MEMBER_PASSWORD)
    if member_role != "member":
        fail(f"member login role mismatch: {member_role}")
    print("OK: member login")

    _, member_writable = request(
        "GET",
        "/api/log-heads/writable",
        token=member_token,
        expected_status=200,
    )
    if isinstance(member_writable, list):
        if not any(get_field(head, "id") == log_head_id for head in member_writable):
            fail("member writable list missing created log head")
    print("OK: list writable log heads (member)")

    _, log_content = request(
        "POST",
        "/api/log-contents",
        token=member_token,
        json_body={
            "log_head_id": log_head_id,
            "content": "Test content",
            "date": to_rfc3339_z(now),
        },
        expected_status=201,
    )
    if not isinstance(log_content, dict) or log_content.get("log_head_id") != log_head_id:
        if get_field(log_content, "log_head_id") != log_head_id:
            fail("log content create response invalid")
    print("OK: create log content")

    request("GET", "/api/admin/log-heads", token=admin_token, expected_status=200)
    print("OK: list admin log heads")

    if not KEEP_DATA:
        request(
            "DELETE",
            f"/api/admin/log-heads/{log_head_id}",
            token=admin_token,
            expected_status=204,
        )
        print("OK: delete log head")

        request(
            "DELETE",
            f"/api/admin/accounts/{member_id}",
            token=admin_token,
            expected_status=204,
        )
        print("OK: delete member account")
    else:
        print("KEEP_DATA set; skipping deletes")

    print("All endpoint checks passed.")


if __name__ == "__main__":
    main()
