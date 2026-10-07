"use client";

import { DownloadOutlined, PlusOutlined, ReloadOutlined } from "@ant-design/icons";
import { ProTable, type ProColumns } from "@ant-design/pro-components";
import { App, Button, Card, InputNumber, Modal, Space, Tag, Typography } from "antd";
import dayjs from "dayjs";
import { useCallback, useEffect, useState } from "react";

import { adminDeleteInviteCodes, adminExportInviteCodesUrl, adminGenerateInviteCodes, adminListInviteCodes, type AdminInviteCode } from "@/services/api/admin";
import { useUserStore } from "@/stores/use-user-store";

type StatusFilter = "all" | "unused" | "used";

const statusLabels: Record<string, string> = { unused: "未使用", used: "已使用" };

export default function AdminInviteCodesPage() {
    const { message, modal } = App.useApp();
    const token = useUserStore((state) => state.token);
    const [items, setItems] = useState<AdminInviteCode[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [status, setStatus] = useState<StatusFilter>("all");
    const [isLoading, setIsLoading] = useState(false);
    const [selectedCodes, setSelectedCodes] = useState<string[]>([]);
    const [generateOpen, setGenerateOpen] = useState(false);
    const [generateCount, setGenerateCount] = useState(10);
    const [generating, setGenerating] = useState(false);

    const loadList = useCallback(async (targetPage = page, targetSize = pageSize, targetStatus = status) => {
        if (!token) return;
        setIsLoading(true);
        try {
            const result = await adminListInviteCodes(token, { page: targetPage, pageSize: targetSize, status: targetStatus === "all" ? undefined : targetStatus });
            setItems(result.items || []);
            setTotal(result.total || 0);
        } catch (error) {
            message.error(error instanceof Error ? error.message : "加载邀请码失败");
        } finally {
            setIsLoading(false);
        }
    }, [message, page, pageSize, status, token]);

    useEffect(() => {
        void loadList();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [page, pageSize, status]);

    const generate = async () => {
        setGenerating(true);
        try {
            const codes = await adminGenerateInviteCodes(token, { count: generateCount });
            setGenerateOpen(false);
            modal.success({
                title: `已生成 ${codes.length} 个邀请码`,
                width: 520,
                content: (
                    <div style={{ maxHeight: 320, overflowY: "auto", fontFamily: "monospace", fontSize: 13, lineHeight: 1.8 }}>
                        {codes.map((code) => <div key={code}>{code}</div>)}
                    </div>
                ),
            });
            await loadList(1, pageSize, status);
        } catch (error) {
            message.error(error instanceof Error ? error.message : "生成邀请码失败");
        } finally {
            setGenerating(false);
        }
    };

    const exportUnused = async () => {
        try {
            const response = await fetch(adminExportInviteCodesUrl(), { headers: { Authorization: `Bearer ${token}` } });
            if (!response.ok) throw new Error(await response.text());
            const text = await response.text();
            const blob = new Blob([text], { type: "text/plain;charset=utf-8" });
            const url = URL.createObjectURL(blob);
            const anchor = document.createElement("a");
            anchor.href = url;
            anchor.download = `invite-codes-${dayjs().format("YYYYMMDD-HHmmss")}.txt`;
            anchor.click();
            URL.revokeObjectURL(url);
        } catch (error) {
            message.error(error instanceof Error ? error.message : "导出失败");
        }
    };

    const deleteSelected = async () => {
        if (!selectedCodes.length) return;
        await adminDeleteInviteCodes(token, { codes: selectedCodes });
        message.success("已删除所选邀请码");
        setSelectedCodes([]);
        await loadList();
    };

    const columns: ProColumns<AdminInviteCode>[] = [
        { title: "邀请码", dataIndex: "code", render: (_, item) => <Typography.Text copyable style={{ fontFamily: "monospace" }}>{item.code}</Typography.Text> },
        { title: "状态", dataIndex: "status", width: 110, render: (_, item) => <Tag color={item.status === "unused" ? "green" : "default"}>{statusLabels[item.status] || item.status}</Tag> },
        { title: "使用者", dataIndex: "usedByName", width: 160, render: (_, item) => item.usedByName || "-" },
        { title: "使用时间", dataIndex: "usedAt", width: 180, render: (_, item) => item.usedAt ? dayjs(item.usedAt).format("YYYY-MM-DD HH:mm:ss") : "-" },
        { title: "创建时间", dataIndex: "createdAt", width: 180, render: (_, item) => item.createdAt ? dayjs(item.createdAt).format("YYYY-MM-DD HH:mm:ss") : "-" },
    ];

    return (
        <Card variant="borderless" title="邀请码管理" extra={
            <Space>
                <Button icon={<DownloadOutlined />} onClick={() => void exportUnused()}>导出未使用</Button>
                <Button danger disabled={!selectedCodes.length} onClick={() => void deleteSelected()}>删除所选</Button>
                <Button icon={<ReloadOutlined />} onClick={() => void loadList()} />
                <Button type="primary" icon={<PlusOutlined />} onClick={() => setGenerateOpen(true)}>生成邀请码</Button>
            </Space>
        }>
            <ProTable<AdminInviteCode>
                rowKey="code"
                search={false}
                options={false}
                loading={isLoading}
                dataSource={items}
                columns={columns}
                rowSelection={{ selectedRowKeys: selectedCodes, onChange: (keys) => setSelectedCodes(keys as string[]) }}
                pagination={{
                    current: page,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    onChange: (nextPage, nextSize) => { setPage(nextPage); setPageSize(nextSize); },
                }}
                headerTitle={
                    <Space>
                        状态
                        <Space.Compact>
                            {(["all", "unused", "used"] as StatusFilter[]).map((value) => (
                                <Button key={value} size="small" type={status === value ? "primary" : "default"} onClick={() => { setPage(1); setStatus(value); }}>
                                    {value === "all" ? "全部" : statusLabels[value]}
                                </Button>
                            ))}
                        </Space.Compact>
                    </Space>
                }
            />
            <Modal
                title="生成邀请码"
                open={generateOpen}
                onCancel={() => setGenerateOpen(false)}
                onOk={() => void generate()}
                confirmLoading={generating}
                okText="生成"
            >
                <Space>
                    <Typography.Text>数量</Typography.Text>
                    <InputNumber min={1} max={500} value={generateCount} onChange={(value) => setGenerateCount(Number(value) || 1)} />
                </Space>
            </Modal>
        </Card>
    );
}
