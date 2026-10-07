"use client";

import { GiftOutlined } from "@ant-design/icons";
import { ProTable, type ProColumns } from "@ant-design/pro-components";
import { App, Button, Card, Col, Input, Row, Statistic, Tag, Typography } from "antd";
import dayjs from "dayjs";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

import { fetchMyCreditLogs, redeemCode, type MyCreditLog } from "@/services/api/auth";
import { useUserStore } from "@/stores/use-user-store";

const creditLogTypeLabels: Record<string, string> = {
    admin_adjust: "后台调整",
    ai_consume: "模型消费",
    ai_refund: "失败返还",
    redeem: "兑换码兑换",
};

export default function MyAccountPage() {
    const { message } = App.useApp();
    const router = useRouter();
    const token = useUserStore((state) => state.token);
    const user = useUserStore((state) => state.user);
    const isReady = useUserStore((state) => state.isReady);
    const hydrateUser = useUserStore((state) => state.hydrateUser);

    useEffect(() => {
        if (isReady && !token) router.replace("/login?redirect=/me");
    }, [isReady, router, token]);
    const [logs, setLogs] = useState<MyCreditLog[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [isLoading, setIsLoading] = useState(false);
    const [code, setCode] = useState("");
    const [redeeming, setRedeeming] = useState(false);

    const loadLogs = useCallback(async (targetPage = page, targetSize = pageSize) => {
        if (!token) return;
        setIsLoading(true);
        try {
            const result = await fetchMyCreditLogs(token, { page: targetPage, pageSize: targetSize });
            setLogs(result.items || []);
            setTotal(result.total || 0);
        } catch {
            // 流水加载失败不打断页面
        } finally {
            setIsLoading(false);
        }
    }, [page, pageSize, token]);

    useEffect(() => {
        void loadLogs();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [page, pageSize]);

    const submitRedeem = async () => {
        if (!code.trim()) {
            message.warning("请填写兑换码");
            return;
        }
        setRedeeming(true);
        try {
            const result = await redeemCode(token, { code: code.trim() });
            message.success(`兑换成功，当前算力点 ${result.credits}`);
            setCode("");
            await hydrateUser();
            await loadLogs(1, pageSize);
            setPage(1);
        } catch (error) {
            message.error(error instanceof Error ? error.message : "兑换失败");
        } finally {
            setRedeeming(false);
        }
    };

    const columns: ProColumns<MyCreditLog>[] = [
        { title: "类型", dataIndex: "type", width: 130, render: (_, item) => <Tag>{creditLogTypeLabels[item.type] || item.type || "-"}</Tag> },
        { title: "变动", dataIndex: "amount", width: 110, render: (_, item) => <Typography.Text type={item.amount >= 0 ? "success" : "danger"}>{item.amount}</Typography.Text> },
        { title: "余额", dataIndex: "balance", width: 110 },
        { title: "备注", dataIndex: "remark", ellipsis: true, render: (_, item) => <Typography.Text type="secondary">{item.remark || "-"}</Typography.Text> },
        { title: "时间", dataIndex: "createdAt", width: 180, render: (_, item) => <Typography.Text type="secondary">{item.createdAt ? dayjs(item.createdAt).format("YYYY-MM-DD HH:mm:ss") : "-"}</Typography.Text> },
    ];

    return (
        <main className="mx-auto w-full max-w-5xl px-4 py-8">
            <Typography.Title level={3} style={{ marginBottom: 20 }}>我的账户</Typography.Title>
            <Row gutter={[16, 16]}>
                <Col xs={24} md={8}>
                    <Card variant="borderless">
                        <Statistic title="当前算力点" value={user?.credits ?? 0} precision={2} />
                        <Typography.Text type="secondary">用户名：{user?.username || "-"}</Typography.Text>
                    </Card>
                </Col>
                <Col xs={24} md={16}>
                    <Card variant="borderless" title="兑换码兑换">
                        <Input.Search
                            enterButton={<><GiftOutlined /> 兑换</>}
                            placeholder="输入管理员发放的兑换码"
                            value={code}
                            onChange={(event) => setCode(event.target.value)}
                            onSearch={() => void submitRedeem()}
                            loading={redeeming}
                            size="large"
                        />
                        <Typography.Text type="secondary" style={{ display: "block", marginTop: 8 }}>
                            兑换成功后，算力点立即到账，可在“算力点日志”中查看明细。
                        </Typography.Text>
                    </Card>
                </Col>
            </Row>
            <Card variant="borderless" title="算力点日志" style={{ marginTop: 16 }}>
                <ProTable<MyCreditLog>
                    rowKey="id"
                    search={false}
                    options={false}
                    loading={isLoading}
                    dataSource={logs}
                    columns={columns}
                    pagination={{
                        current: page,
                        pageSize,
                        total,
                        showSizeChanger: true,
                        onChange: (nextPage, nextSize) => { setPage(nextPage); setPageSize(nextSize); },
                    }}
                />
            </Card>
        </main>
    );
}
