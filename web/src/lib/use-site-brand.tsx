"use client";

import { useEffect } from "react";

import { useConfigStore } from "@/stores/use-config-store";

export const DEFAULT_SITE_NAME = "无限画布";
export const DEFAULT_LOGO_URL = "/logo.svg";

export function useSiteBrand() {
    const siteName = useConfigStore((state) => state.publicSettings?.site?.siteName) || DEFAULT_SITE_NAME;
    const logoUrl = useConfigStore((state) => state.publicSettings?.site?.logoUrl) || DEFAULT_LOGO_URL;
    return { siteName, logoUrl };
}

/** 同步浏览器标签页标题为后台配置的站点名称。 */
export function SiteTitle() {
    const { siteName } = useSiteBrand();
    useEffect(() => {
        document.title = siteName;
    }, [siteName]);
    return null;
}
