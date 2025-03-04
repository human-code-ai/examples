import axios, { AxiosInstance, AxiosResponse } from 'axios';
import * as crypto from 'crypto';

interface ClientResponse<T> {
  code: number;
  msg: string;
  result: T;
}

interface GetSessionIdResult {
  session_id: string;
}

interface VerifyResult {
  human_id: string;
}

interface ClientConfig {
  baseUrl: string,
  debug: boolean,
  appId: string,
  appKey: string
}

class HumanCodeProvider {
  private readonly apiClient: AxiosInstance;

  constructor(
    private readonly config: ClientConfig
  ) {
    this.apiClient = axios.create({
      baseURL: config.baseUrl,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (config.debug) {
      this.addDebugInterceptors();
    }
  }

  private addDebugInterceptors() {
    this.apiClient.interceptors.request.use(request => {
      console.log('Request:', JSON.stringify(request, null, 2));
      return request;
    });

    this.apiClient.interceptors.response.use(response => {
      console.log('Response:', JSON.stringify(response.data, null, 2));
      return response;
    });
  }

  private genSign(data: string): string {
    console.debug(`🚀 humancode.provider.ts[56] - data: `, data);
    const hmac = crypto.createHmac('sha256', this.config.appKey);
    hmac.update(data);
    return hmac.digest('hex');
  }

  getConfig(): ClientConfig {
    return this.config;
  }

  async getSessionId(nonceStr: string): Promise<GetSessionIdResult> {
    const timestamp = Date.now();
    const postBody = {
      timestamp: timestamp.toString(),
      nonce_str: nonceStr
    };

    const sign = this.genSign(JSON.stringify(postBody));
    const url = `/api/session/v2/get_id?app_id=${this.config.appId}&sign=${sign}`;

    try {
      const response = await this.apiClient.post<ClientResponse<GetSessionIdResult>>(url, postBody);
      return this.handleResponse(response);
    } catch (error) {
      throw new Error(`Request failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  genRegistrationUrl(sessionId: string, callBackUrl: string): string {
    const timestamp = Date.now();
    return `${this.config.baseUrl}/authentication/index.html?session_id=${sessionId}&callback_url=${callBackUrl}&ts=${timestamp}#/`;
  }

  genVerificationUrl(sessionId: string, humanId: string, callBackUrl: string): string {
    const timestamp = Date.now();
    return `${this.config.baseUrl}/api/authentication/index.html?session_id=${sessionId}&human_id=${humanId}&callback_url=${callBackUrl}&ts=${timestamp}#/`;
  }

  async verify(sessionId: string, vCode: string, nonceStr: string): Promise<VerifyResult> {
    const timestamp = Date.now();
    const postBody = {
      session_id: sessionId,
      vcode: vCode,
      timestamp: timestamp.toString(),
      nonce_str: nonceStr
    };

    const sign = this.genSign(JSON.stringify(postBody));
    const url = `/api/vcode/v2/verify?app_id=${this.config.appId}&sign=${sign}`;

    try {
      const response = await this.apiClient.post<ClientResponse<VerifyResult>>(url, postBody);
      return this.handleResponse(response);
    } catch (error) {
      throw new Error(`Request failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  private handleResponse<T>(response: AxiosResponse<ClientResponse<T>>): T {
    if (response.status !== 200) {
      throw new Error(`HTTP Error: ${response.status}`);
    }

    const { code, msg, result } = response.data;
    if (code !== 0) {
      throw new Error(`API Error: [${code}] ${msg}`);
    }

    return result;
  }
}

export default HumanCodeProvider;