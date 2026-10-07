import { apiGet, apiPost } from "@/services/api/request";

export const AUTH_TOKEN_KEY = "infinite-canvas-auth-token-v1";

export type UserRole = "guest" | "user" | "admin";

export type AuthUser = {
    id: string;
    username: string;
    displayName: string;
    avatarUrl: string;
    role: UserRole;
    credits: number;
    createdAt: string;
    updatedAt: string;
};

export type AuthSession = {
    token: string;
    user: AuthUser;
};

export type AuthPayload = {
    username: string;
    password: string;
    inviteCode?: string;
};

export async function login(payload: AuthPayload) {
    return apiPost<AuthSession>("/api/auth/login", payload);
}

export async function register(payload: AuthPayload) {
    return apiPost<AuthSession>("/api/auth/register", payload);
}

export async function fetchCurrentUser(token?: string) {
    return apiGet<AuthUser>("/api/auth/me", undefined, token);
}


export type MyCreditLog = {
    id: string;
    type: string;
    amount: number;
    balance: number;
    relatedId: string;
    remark: string;
    createdAt: string;
};

export type MyCreditLogList = {
    items: MyCreditLog[];
    total: number;
};

export async function redeemCode(token: string, payload: { code: string }) {
    return apiPost<AuthUser>("/api/v1/redeem", payload, token);
}

export async function fetchMyCreditLogs(token: string, params: { page?: number; pageSize?: number }) {
    return apiGet<MyCreditLogList>("/api/v1/credit-logs", params, token);
}
