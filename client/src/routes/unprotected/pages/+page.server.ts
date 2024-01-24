import { error, json } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({params}) => {
    const r = await fetch("http://resume-service.local/candidate");
    const res = await r.json();

    const c = await fetch("http://candidate-service.local/ping");
    const can = await c.json();

    console.log("ADSF ", res, "\nQWER ", can)

    return {res: res, can: can};
}