import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class FaturamentoService {
  
  private apiUrl = 'http://localhost:8081/faturar';

  constructor(private http: HttpClient) { }

  emitirNota(dadosFaturamento: any): Observable<any> {
    return this.http.post(this.apiUrl, dadosFaturamento);
  }
}