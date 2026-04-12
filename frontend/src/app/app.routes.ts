import { Routes } from '@angular/router';
import { EstoqueComponent } from './estoque/estoque.component';
import { FaturamentoComponent } from './faturamento/faturamento.component';

export const routes: Routes = [
    {path: 'estoque', component: EstoqueComponent},
    {path : 'faturamento', component: FaturamentoComponent},
    {path: '', redirectTo: '/estoque', pathMatch: 'full'}
];
