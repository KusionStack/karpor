/* eslint-disable */
// 该文件由 OneAPI 自动生成，请勿手动修改！
import { request } from '@alipay/bigfish';

/** api_v1_dw_usingPOST 查询计算机流列表 POST /api/v1/dw */
export async function apiV1DwUsingPOST(
  body?: {
    name?: string;
  },
  options?: { [key: string]: any },
) {
  return request<{
    success?: string;
    errorMessage?: string;
    data: {
      current?: string;
      pageSize?: string;
      total?: string;
      list?: Array<{
        id?: string;
        name?: string;
        sourceCode?: string;
        gmtModified?: string;
        code?: string;
        operator?: string;
      }>;
    };
  }>('/api/v1/dw', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** api_v1_file_queryAlgorithmFile_usingGET 查询算法文件 GET /api/v1/file/queryAlgorithmFile */
export async function apiV1FileQueryAlgorithmFileUsingGET(options?: {
  [key: string]: any;
}) {
  return request<{ success?: boolean; data: Record<string, any> }>(
    '/api/v1/file/queryAlgorithmFile',
    {
      method: 'GET',
      ...(options || {}),
    },
  );
}

/** 此处后端没有提供注释 GET /api/queryRecStrategy */
export async function queryRecStrategy(options?: { [key: string]: any }) {
  return request<{
    success?: boolean;
    resultCode?: string;
    resultDesc?: string;
    resultView?: string;
    result?: Array<{
      id?: number;
      name?: string;
      stickList?: Array<Record<string, any>>;
      crowdId?: any;
      crowdDesc?: any;
      recall?: string;
      type?: string;
      blackList?: Array<Record<string, any>>;
      interfaceType?: string;
      gmtCreate?: string;
      gmtModified?: string;
      recallStrategyDTO?: any;
      interfaceModeParam: Record<string, any>;
      strategyRecallInfo: {
        exception?: string;
        took?: number;
        total?: number;
        cost?: number;
        size?: number;
        dsl?: string;
        content?: string;
      };
    }>;
    externalData: Record<string, any>;
  }>('/api/queryRecStrategy', {
    method: 'GET',
    ...(options || {}),
  });
}
